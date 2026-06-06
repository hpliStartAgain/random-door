from __future__ import annotations

import binascii
import html
import json
import re
import struct
import time
import zlib
from pathlib import Path
from urllib.parse import quote

import requests

try:
    from PIL import Image, ImageDraw, ImageFont
except ImportError:  # pragma: no cover - only used when requirements are missing.
    Image = None
    ImageDraw = None
    ImageFont = None


COMMONS_API = "https://commons.wikimedia.org/w/api.php"
WIKIDATA_API = "https://www.wikidata.org/w/api.php"
ZH_WIKIPEDIA_API = "https://zh.wikipedia.org/w/api.php"
EN_WIKIPEDIA_API = "https://en.wikipedia.org/w/api.php"
USER_AGENT = "RandomDoorSeedBuilder/0.1 (local demo asset pipeline)"
BAD_TITLE_PARTS = (
    "map",
    "locator",
    "blank",
    "flag",
    "seal",
    "logo",
    "icon",
    "svg",
    "diagram",
    "route",
    "map of",
    "political",
    "emblem",
    "coat of arms",
)


def collect_static_urls(cities: list[dict], achievements: list[dict]) -> list[str]:
    urls: list[str] = []
    for city in cities:
        urls.append(city["cover_image_url"])
        for landmark in city.get("landmarks", []):
            urls.append(landmark["image_url"])
        for food in city.get("foods", []):
            urls.append(food["image_url"])
        for character in city.get("characters", []):
            urls.append(character["avatar_url"])
    for achievement in achievements:
        urls.append(achievement["badge_url"])
    return sorted(set(urls))


def ensure_assets(repo_root: Path, cities: list[dict], achievements: list[dict], offline: bool = False) -> list[dict]:
    manifest: list[dict] = []
    session = requests.Session()
    session.headers.update({"User-Agent": USER_AGENT})

    for spec in iter_asset_specs(cities, achievements):
        url = spec["url"]
        if not url.startswith("/static/"):
            raise ValueError(f"asset URL must start with /static/: {url}")

        target = repo_root / "backend" / url.lstrip("/")
        target.parent.mkdir(parents=True, exist_ok=True)
        width, height = dimensions_for(url)

        entry: dict
        if not offline and spec["kind"] in {"cover", "landmark", "food"}:
            entry = fetch_commons_asset(session, target, spec, width, height)
        else:
            entry = None

        if entry is None:
            seed = color_seed(url)
            if Image is not None:
                write_label_image(target, width, height, seed, spec["label"], spec["kind"])
            else:
                write_png(target, width, height, seed)
            entry = {
                "source": "local project-generated illustration",
                "license": "project-generated",
                "source_url": None,
            }

        entry.update(
            {
                "url": url,
                "path": str(target.relative_to(repo_root)).replace("\\", "/"),
                "kind": spec["kind"],
                "label": spec["label"],
                "width": width,
                "height": height,
            }
        )
        manifest.append(entry)
    return manifest


def iter_asset_specs(cities: list[dict], achievements: list[dict]) -> list[dict]:
    specs: list[dict] = []
    for city in cities:
        landmarks = city.get("landmarks", [])
        cover_query = f"{city['name']} {landmarks[0]['name'] if landmarks else ''} 城市 风景"
        specs.append(
            {
                "url": city["cover_image_url"],
                "kind": "cover",
                "label": f"{city['name']}城市封面",
                "query": cover_query,
                "wikidata": city["name"],
            }
        )
        for landmark in landmarks:
            specs.append(
                {
                    "url": landmark["image_url"],
                    "kind": "landmark",
                    "label": f"{city['name']} {landmark['name']}",
                    "query": f"{landmark['name']} {city['name']}",
                    "wikidata": landmark["name"],
                }
            )
        for food in city.get("foods", []):
            specs.append(
                {
                    "url": food["image_url"],
                    "kind": "food",
                    "label": food["name"],
                    "query": f"{food['name']} 中国 美食",
                    "wikidata": food["name"],
                }
            )
        for character in city.get("characters", []):
            specs.append(
                {
                    "url": character["avatar_url"],
                    "kind": "character",
                    "label": character["name"],
                    "query": "",
                }
            )
    for achievement in achievements:
        specs.append(
            {
                "url": achievement["badge_url"],
                "kind": "badge",
                "label": achievement["name"],
                "query": "",
            }
        )
    return specs


def fetch_commons_asset(
    session: requests.Session,
    target: Path,
    spec: dict,
    width: int,
    height: int,
) -> dict | None:
    wikipedia_candidate = image_from_wikipedia(session, spec.get("wikidata", ""), width)
    if wikipedia_candidate is not None:
        try:
            image_url = wikipedia_candidate.get("thumburl") or wikipedia_candidate.get("url")
            if image_url:
                response = session.get(image_url, timeout=12)
                response.raise_for_status()
                save_image_bytes(target, response.content, width, height)
                return {
                    "source": "Wikimedia project page image",
                    "license": "see source page image metadata",
                    "source_url": wikipedia_candidate.get("source_page"),
                    "author": None,
                    "title": wikipedia_candidate.get("title"),
                    "resolver": wikipedia_candidate.get("resolver"),
                }
        except Exception:
            pass

    wikidata_candidate = image_from_wikidata(session, spec.get("wikidata", ""), width)
    if wikidata_candidate is not None:
        try:
            image_url = wikidata_candidate.get("thumburl") or wikidata_candidate.get("url")
            if image_url:
                response = session.get(image_url, timeout=12)
                response.raise_for_status()
                save_image_bytes(target, response.content, width, height)
                return {
                    "source": "Wikimedia Commons",
                    "license": wikidata_candidate.get("license") or "see Wikimedia Commons source page",
                    "source_url": wikidata_candidate.get("descriptionurl"),
                    "author": wikidata_candidate.get("artist"),
                    "title": wikidata_candidate.get("title"),
                    "resolver": "wikidata_p18",
                }
        except Exception:
            pass

    candidates = search_commons(session, spec["query"], width)
    for candidate in candidates:
        try:
            image_url = candidate.get("thumburl") or candidate.get("url")
            if not image_url:
                continue
            response = session.get(image_url, timeout=12)
            response.raise_for_status()
            save_image_bytes(target, response.content, width, height)
            return {
                "source": "Wikimedia Commons",
                "license": candidate.get("license") or "see Wikimedia Commons source page",
                "source_url": candidate.get("descriptionurl"),
                "author": candidate.get("artist"),
                "title": candidate.get("title"),
                "resolver": "commons_search",
            }
        except Exception:
            continue
    return None


def image_from_wikipedia(session: requests.Session, term: str, thumb_width: int) -> dict | None:
    for api, lang in ((ZH_WIKIPEDIA_API, "zhwiki"),):
        for title in wikipedia_title_variants(term):
            try:
                response = session.get(
                    api,
                    params={
                        "action": "query",
                        "format": "json",
                        "prop": "pageimages|info",
                        "piprop": "original|thumbnail",
                        "pithumbsize": min(max(thumb_width, 512), 1600),
                        "inprop": "url",
                        "redirects": 1,
                        "titles": title,
                    },
                    timeout=8,
                )
                response.raise_for_status()
                pages = response.json().get("query", {}).get("pages", {})
            except Exception:
                continue
            for page in pages.values():
                if page.get("missing") is not None:
                    continue
                thumbnail = page.get("thumbnail", {}).get("source")
                original = page.get("original", {}).get("source")
                image_url = thumbnail or original
                if not image_url:
                    continue
                if not image_url.startswith("https://upload.wikimedia.org/"):
                    continue
                return {
                    "title": page.get("title", title),
                    "url": original or thumbnail,
                    "thumburl": thumbnail,
                    "source_page": page.get("fullurl"),
                    "resolver": f"{lang}_pageimage",
                }
    return None


def wikipedia_title_variants(term: str) -> list[str]:
    term = term.strip()
    if not term:
        return []
    variants = [term]
    match = re.search(r"(.+)[（(](.+)[）)]", term)
    if match:
        outside = match.group(1).strip()
        inside = match.group(2).strip()
        variants.extend([inside + outside, outside, inside])
    if term.endswith("城市封面"):
        variants.append(term.removesuffix("城市封面"))
    return dedupe(variants)


def dedupe(values: list[str]) -> list[str]:
    seen = set()
    result = []
    for value in values:
        if value and value not in seen:
            seen.add(value)
            result.append(value)
    return result


def image_from_wikidata(session: requests.Session, term: str, thumb_width: int) -> dict | None:
    term = term.strip()
    if not term:
        return None
    try:
        search_resp = session.get(
            WIKIDATA_API,
            params={
                "action": "wbsearchentities",
                "format": "json",
                "language": "zh",
                "uselang": "zh",
                "search": term,
                "limit": 5,
            },
            timeout=8,
        )
        search_resp.raise_for_status()
        search_data = search_resp.json()
    except Exception:
        return None

    for hit in search_data.get("search", []):
        qid = hit.get("id")
        if not qid:
            continue
        filename = p18_filename(session, qid)
        if not filename:
            continue
        info = commons_file_info(session, filename, thumb_width)
        if info is not None:
            return info
    return None


def p18_filename(session: requests.Session, qid: str) -> str | None:
    try:
        entity_resp = session.get(
            WIKIDATA_API,
            params={
                "action": "wbgetentities",
                "format": "json",
                "ids": qid,
                "props": "claims",
            },
            timeout=8,
        )
        entity_resp.raise_for_status()
        entity = entity_resp.json().get("entities", {}).get(qid, {})
        claims = entity.get("claims", {}).get("P18") or []
    except Exception:
        return None
    for claim in claims:
        value = (
            claim.get("mainsnak", {})
            .get("datavalue", {})
            .get("value")
        )
        if isinstance(value, str) and value:
            return value
    return None


def commons_file_info(session: requests.Session, filename: str, thumb_width: int) -> dict | None:
    title = "File:" + filename
    if bad_commons_title(title):
        return None
    try:
        response = session.get(
            COMMONS_API,
            params={
                "action": "query",
                "format": "json",
                "titles": title,
                "prop": "imageinfo",
                "iiprop": "url|mime|size|extmetadata",
                "iiurlwidth": min(max(thumb_width, 512), 1600),
            },
            timeout=8,
        )
        response.raise_for_status()
        pages = response.json().get("query", {}).get("pages", {})
    except Exception:
        return None
    for page in pages.values():
        imageinfo = (page.get("imageinfo") or [{}])[0]
        mime = imageinfo.get("mime", "")
        if mime not in {"image/jpeg", "image/png", "image/webp"}:
            return None
        meta = imageinfo.get("extmetadata") or {}
        return {
            "title": page.get("title", title),
            "url": imageinfo.get("url"),
            "thumburl": imageinfo.get("thumburl"),
            "descriptionurl": imageinfo.get("descriptionurl"),
            "license": metadata_value(meta, "LicenseShortName") or metadata_value(meta, "UsageTerms"),
            "artist": clean_metadata(metadata_value(meta, "Artist")),
        }
    return None


def search_commons(session: requests.Session, query: str, thumb_width: int) -> list[dict]:
    search_terms = [query, query.replace(" 中国 美食", ""), query.replace(" 城市 风景", "")]
    pages: list[dict] = []
    seen_titles: set[str] = set()
    for term in search_terms:
        term = re.sub(r"\s+", " ", term).strip()
        if not term:
            continue
        params = {
            "action": "query",
            "format": "json",
            "generator": "search",
            "gsrnamespace": 6,
            "gsrsearch": term,
            "gsrlimit": 8,
            "prop": "imageinfo",
            "iiprop": "url|mime|size|extmetadata",
            "iiurlwidth": min(max(thumb_width, 512), 1600),
        }
        try:
            response = session.get(COMMONS_API, params=params, timeout=8)
            response.raise_for_status()
            data = response.json()
        except Exception:
            continue
        for page in data.get("query", {}).get("pages", {}).values():
            title = page.get("title", "")
            if title in seen_titles or bad_commons_title(title):
                continue
            imageinfo = (page.get("imageinfo") or [{}])[0]
            mime = imageinfo.get("mime", "")
            if mime not in {"image/jpeg", "image/png", "image/webp"}:
                continue
            meta = imageinfo.get("extmetadata") or {}
            pages.append(
                {
                    "title": title,
                    "url": imageinfo.get("url"),
                    "thumburl": imageinfo.get("thumburl"),
                    "descriptionurl": imageinfo.get("descriptionurl"),
                    "license": metadata_value(meta, "LicenseShortName") or metadata_value(meta, "UsageTerms"),
                    "artist": clean_metadata(metadata_value(meta, "Artist")),
                    "width": imageinfo.get("width", 0),
                    "height": imageinfo.get("height", 0),
                    "score": candidate_score(title, term, imageinfo),
                }
            )
            seen_titles.add(title)
        time.sleep(0.05)
    return sorted(pages, key=lambda item: item["score"], reverse=True)


def bad_commons_title(title: str) -> bool:
    lower = title.lower()
    return any(part in lower for part in BAD_TITLE_PARTS)


def candidate_score(title: str, query: str, imageinfo: dict) -> int:
    score = 0
    lower_title = title.lower()
    for token in re.split(r"\s+", query.lower()):
        token = token.strip()
        if token and token in lower_title:
            score += 5
    width = int(imageinfo.get("width") or 0)
    height = int(imageinfo.get("height") or 0)
    if width >= 900 and height >= 600:
        score += 6
    if width > height:
        score += 3
    if "quality images" in lower_title or "featured" in lower_title:
        score += 2
    return score


def metadata_value(meta: dict, key: str) -> str | None:
    raw = meta.get(key)
    if not raw:
        return None
    return raw.get("value")


def clean_metadata(value: str | None) -> str | None:
    if not value:
        return None
    value = re.sub(r"<[^>]+>", "", value)
    return html.unescape(value).strip() or None


def save_image_bytes(target: Path, data: bytes, width: int, height: int) -> None:
    if Image is None:
        target.write_bytes(data)
        return
    from io import BytesIO

    with Image.open(BytesIO(data)) as source:
        image = source.convert("RGB")
        image = crop_to_aspect(image, width / height)
        image = image.resize((width, height), Image.Resampling.LANCZOS)
        image.save(target, format="PNG", optimize=True)


def crop_to_aspect(image, target_aspect: float):
    width, height = image.size
    current = width / height
    if current > target_aspect:
        new_width = int(height * target_aspect)
        left = (width - new_width) // 2
        return image.crop((left, 0, left + new_width, height))
    if current < target_aspect:
        new_height = int(width / target_aspect)
        top = (height - new_height) // 2
        return image.crop((0, top, width, top + new_height))
    return image


def dimensions_for(url: str) -> tuple[int, int]:
    if "/characters/" in url or "/badges/" in url:
        return 512, 512
    if "/foods/" in url:
        return 1024, 768
    return 1280, 720


def color_seed(value: str) -> tuple[int, int, int]:
    checksum = binascii.crc32(value.encode("utf-8"))
    r = 80 + checksum % 140
    g = 80 + (checksum >> 8) % 140
    b = 80 + (checksum >> 16) % 140
    return r, g, b


def write_png(path: Path, width: int, height: int, seed: tuple[int, int, int]) -> None:
    r0, g0, b0 = seed
    rows = bytearray()
    for y in range(height):
        rows.append(0)
        for x in range(width):
            shade = (x * 31 // max(width - 1, 1) + y * 47 // max(height - 1, 1)) % 96
            rows.extend(
                (
                    (r0 + shade) % 256,
                    (g0 + shade // 2) % 256,
                    (b0 + shade // 3) % 256,
                )
            )

    def chunk(kind: bytes, data: bytes) -> bytes:
        return (
            struct.pack(">I", len(data))
            + kind
            + data
            + struct.pack(">I", binascii.crc32(kind + data) & 0xFFFFFFFF)
        )

    header = struct.pack(">IIBBBBB", width, height, 8, 2, 0, 0, 0)
    payload = b"\x89PNG\r\n\x1a\n"
    payload += chunk(b"IHDR", header)
    payload += chunk(b"IDAT", zlib.compress(bytes(rows), level=9))
    payload += chunk(b"IEND", b"")
    path.write_bytes(payload)


def write_label_image(path: Path, width: int, height: int, seed: tuple[int, int, int], label: str, kind: str) -> None:
    if Image is None:
        write_png(path, width, height, seed)
        return
    r0, g0, b0 = seed
    image = Image.new("RGB", (width, height), (r0, g0, b0))
    pixels = image.load()
    for y in range(height):
        for x in range(width):
            shade = (x * 29 // max(width - 1, 1) + y * 41 // max(height - 1, 1)) % 80
            pixels[x, y] = ((r0 + shade) % 256, (g0 + shade // 2) % 256, (b0 + shade // 3) % 256)

    draw = ImageDraw.Draw(image)
    font = default_font(42 if width >= 1000 else 30)
    small_font = default_font(24 if width >= 1000 else 18)
    text = label[:18]
    subtitle = {"character": "AI 人物插画", "badge": "成就徽章"}.get(kind, "项目生成素材")
    draw.rectangle((0, int(height * 0.68), width, height), fill=(0, 0, 0))
    draw.text((32, int(height * 0.72)), text, fill=(255, 255, 255), font=font)
    draw.text((32, int(height * 0.72) + 56), subtitle, fill=(220, 220, 220), font=small_font)
    image.save(path, format="PNG", optimize=True)


def default_font(size: int):
    if ImageFont is None:
        return None
    for name in ("msyh.ttc", "simhei.ttf", "arial.ttf"):
        try:
            return ImageFont.truetype(name, size)
        except Exception:
            continue
    return ImageFont.load_default()


def write_manifest(repo_root: Path, manifest: list[dict]) -> None:
    target = repo_root / "scripts" / "asset_manifest.json"
    target.write_text(json.dumps(manifest, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")


def assert_static_files_exist(repo_root: Path, cities: list[dict], achievements: list[dict]) -> None:
    missing = []
    for url in collect_static_urls(cities, achievements):
        path = repo_root / "backend" / url.lstrip("/")
        if not path.is_file():
            missing.append(url)
    if missing:
        raise FileNotFoundError("missing static assets:\n" + "\n".join(missing))
