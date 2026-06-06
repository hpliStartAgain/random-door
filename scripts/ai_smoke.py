from __future__ import annotations

import argparse
import base64
import json
import mimetypes
import os
import sys
from pathlib import Path
from typing import Any

import requests


REPO_ROOT = Path(__file__).resolve().parents[1]
DEFAULT_ENV = REPO_ROOT / ".env"


def load_env(path: Path) -> None:
    if not path.exists():
        print(f"env: missing {path}")
        return
    for raw in path.read_text(encoding="utf-8").splitlines():
        line = raw.strip()
        if not line or line.startswith("#") or "=" not in line:
            continue
        key, value = line.split("=", 1)
        key = key.strip()
        value = value.strip().strip('"').strip("'")
        if key and key not in os.environ:
            os.environ[key] = value
    print(f"env: loaded {path}")


def configured(name: str) -> bool:
    return bool(os.getenv(name, "").strip())


def require_env(*names: str) -> bool:
    missing = [name for name in names if not configured(name)]
    if missing:
        print(f"skip: missing {', '.join(missing)}")
        return False
    return True


def looks_placeholder(value: str) -> bool:
    normalized = value.lower()
    return any(marker in normalized for marker in ("replace_me", "your-", "your_", "change_me"))


def reject_placeholder(label: str, value: str) -> bool:
    if looks_placeholder(value):
        print(f"skip: {label} still looks like a placeholder")
        return True
    return False


def post_json(url: str, api_key: str, payload: dict[str, Any], timeout: int) -> requests.Response:
    return requests.post(
        url,
        headers={
            "Authorization": f"Bearer {api_key}",
            "Content-Type": "application/json",
        },
        json=payload,
        timeout=timeout,
    )


def safe_error(exc: BaseException) -> str:
    text = str(exc)
    for secret_name in ("LLM_API_KEY", "IMAGE_API_KEY"):
        secret = os.getenv(secret_name, "")
        if secret:
            text = text.replace(secret, "[REDACTED]")
    text = " ".join(text.split())
    if len(text) > 240:
        text = text[:237] + "..."
    return text


def smoke_llm(timeout: int) -> bool:
    if not require_env("LLM_API_BASE", "LLM_API_KEY", "LLM_MODEL"):
        return False
    if reject_placeholder("LLM_API_BASE", os.environ["LLM_API_BASE"]):
        return False
    if reject_placeholder("LLM_API_KEY", os.environ["LLM_API_KEY"]):
        return False

    base = os.environ["LLM_API_BASE"].rstrip("/")
    url = f"{base}/chat/completions"
    payload = {
        "model": os.environ["LLM_MODEL"],
        "messages": [
            {
                "role": "system",
                "content": "你是任意门项目的接口联调助手。用一句中文回答，不要输出密钥。",
            },
            {"role": "user", "content": "请回复：联调成功"},
        ],
        "temperature": 0.2,
    }

    try:
        resp = post_json(url, os.environ["LLM_API_KEY"], payload, timeout)
    except requests.RequestException as exc:
        print(f"llm: request_error type={type(exc).__name__} detail={safe_error(exc)}")
        return False

    if resp.status_code != 200:
        reason = "auth_failed_or_invalid_key" if resp.status_code == 401 else "upstream_error"
        print(f"llm: failed status={resp.status_code} reason={reason} body_len={len(resp.text)}")
        return False

    try:
        data = resp.json()
        reply = data["choices"][0]["message"]["content"]
    except (KeyError, IndexError, TypeError, json.JSONDecodeError):
        print(f"llm: failed status=200 parse_error body_len={len(resp.text)}")
        return False

    print(f"llm: ok status=200 reply_len={len(reply)}")
    return True


def image_data_url(path: Path) -> str:
    mime_type = mimetypes.guess_type(path.name)[0] or "image/png"
    encoded = base64.b64encode(path.read_bytes()).decode("ascii")
    return f"data:{mime_type};base64,{encoded}"


def dashscope_endpoint(base: str) -> str:
    base = base.rstrip("/")
    if base.endswith("/generation"):
        return base
    return f"{base}/api/v1/services/aigc/multimodal-generation/generation"


def find_first_image(value: Any) -> str | None:
    if isinstance(value, dict):
        for key in ("image", "image_url", "url", "b64_json"):
            raw = value.get(key)
            if isinstance(raw, str) and raw:
                return raw
        for child in value.values():
            found = find_first_image(child)
            if found:
                return found
    if isinstance(value, list):
        for child in value:
            found = find_first_image(child)
            if found:
                return found
    return None


def smoke_image(args: argparse.Namespace) -> bool:
    if not require_env("IMAGE_API_BASE", "IMAGE_API_KEY"):
        return False
    if reject_placeholder("IMAGE_API_BASE", os.environ["IMAGE_API_BASE"]):
        return False
    if reject_placeholder("IMAGE_API_KEY", os.environ["IMAGE_API_KEY"]):
        return False
    if os.environ["IMAGE_API_BASE"].strip().lower() == "mock":
        print("image: skipped; IMAGE_API_BASE is mock")
        return False
    if not args.confirm_image_cost:
        print("image: configured; pass --confirm-image-cost to call the paid image API")
        return False
    if not args.selfie:
        print("image: skipped; --selfie is required for real image smoke")
        return False

    selfie = Path(args.selfie)
    if not selfie.exists():
        print(f"image: skipped; selfie not found: {selfie}")
        return False

    ref = Path(args.ref) if args.ref else REPO_ROOT / "backend" / "static" / "landmarks" / "beijing_cover.png"
    if not ref.exists():
        print(f"image: skipped; ref image not found: {ref}")
        return False

    model = os.getenv("IMAGE_MODEL", "wan2.7-image-pro")
    base = os.environ["IMAGE_API_BASE"]
    url = dashscope_endpoint(base)
    payload = {
        "model": model,
        "input": {
            "messages": [
                {
                    "role": "user",
                    "content": [
                        {
                            "text": (
                                "生成一张真实旅行打卡照：保留用户上传自拍中的人物身份，"
                                "自然合成到北京城市地标场景。禁止色情、暴力、侮辱内容。"
                            )
                        },
                        {"image": image_data_url(selfie)},
                        {"image": image_data_url(ref)},
                    ],
                }
            ]
        },
        "parameters": {
            "size": "2K",
            "n": 1,
            "watermark": False,
            "thinking_mode": True,
        },
    }

    try:
        resp = post_json(url, os.environ["IMAGE_API_KEY"], payload, args.timeout)
    except requests.RequestException as exc:
        print(f"image: request_error type={type(exc).__name__} detail={safe_error(exc)}")
        return False

    if resp.status_code != 200:
        reason = "auth_failed_or_invalid_key" if resp.status_code == 401 else "upstream_error"
        print(f"image: failed status={resp.status_code} reason={reason} body_len={len(resp.text)}")
        return False

    try:
        data = resp.json()
    except json.JSONDecodeError:
        print(f"image: failed status=200 parse_error body_len={len(resp.text)}")
        return False

    image_value = find_first_image(data)
    if not image_value:
        print("image: failed status=200 no_image_in_response")
        return False

    print(f"image: ok status=200 result_kind={'url' if image_value.startswith('http') else 'inline'}")
    return True


def main() -> int:
    parser = argparse.ArgumentParser(description="Smoke test external AI endpoints without printing secrets.")
    parser.add_argument("--env", default=str(DEFAULT_ENV), help="env file to load, default .env")
    parser.add_argument("--llm", action="store_true", help="test LLM chat/completions")
    parser.add_argument("--image", action="store_true", help="test DashScope image generation")
    parser.add_argument("--selfie", help="selfie image path for --image")
    parser.add_argument("--ref", help="reference image path for --image")
    parser.add_argument("--confirm-image-cost", action="store_true", help="allow one real image generation request")
    parser.add_argument("--timeout", type=int, default=60, help="request timeout seconds")
    args = parser.parse_args()

    if not args.llm and not args.image:
        args.llm = True

    load_env(Path(args.env))

    ok = True
    if args.llm:
        ok = smoke_llm(args.timeout) and ok
    if args.image:
        ok = smoke_image(args) and ok

    return 0 if ok else 1


if __name__ == "__main__":
    sys.exit(main())
