from __future__ import annotations

import json
import sys
import time
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))

from seed_builder.assets import ensure_assets, write_manifest


REPO_ROOT = Path(__file__).resolve().parents[1]
CITIES_JSON = REPO_ROOT / "backend" / "data" / "seed" / "cities.json"
ACHIEVEMENTS_JSON = REPO_ROOT / "backend" / "data" / "seed" / "achievements.json"


def load_json(path: Path) -> list[dict]:
    return json.loads(path.read_text(encoding="utf-8"))


def main() -> None:
    cities = load_json(CITIES_JSON)
    achievements = load_json(ACHIEVEMENTS_JSON)
    manifest: list[dict] = []
    for index, city in enumerate(cities, start=1):
        manifest.extend(ensure_assets(REPO_ROOT, [city], [], offline=False))
        write_manifest(REPO_ROOT, manifest)
        sourced = sum(1 for item in manifest if item["source"] != "local project-generated illustration")
        print(f"{index:02d}/{len(cities)} {city['name']} sourced={sourced}", flush=True)
        time.sleep(0.2)
    manifest.extend(ensure_assets(REPO_ROOT, [], achievements, offline=True))
    write_manifest(REPO_ROOT, manifest)
    print(f"done total={len(manifest)} sourced={sum(1 for item in manifest if item['source'] != 'local project-generated illustration')}", flush=True)


if __name__ == "__main__":
    main()
