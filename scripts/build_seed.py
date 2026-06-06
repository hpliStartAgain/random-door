from __future__ import annotations

import argparse
import csv
import json
import os
import subprocess
import sys
import time
from pathlib import Path

from seed_builder.assets import assert_static_files_exist, ensure_assets, write_manifest
from seed_builder.catalog import build_catalog


REPO_ROOT = Path(__file__).resolve().parents[1]
SEED_DIR = REPO_ROOT / "backend" / "data" / "seed"
CITIES_JSON = SEED_DIR / "cities.json"
ACHIEVEMENTS_JSON = SEED_DIR / "achievements.json"
SEED_INPUTS = REPO_ROOT / "scripts" / "seed_inputs.csv"


def load_json(path: Path) -> list[dict]:
    return json.loads(path.read_text(encoding="utf-8"))


def write_json(path: Path, data: list[dict]) -> None:
    path.write_text(json.dumps(data, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")


def audit_db(_: argparse.Namespace) -> None:
    rows = []
    mysql_cmd = os.getenv("MYSQL_CLI", "mysql")
    db_name = os.getenv("DB_NAME")
    if db_name:
        query = "SELECT name, province, lat, lng FROM cities ORDER BY id LIMIT 35;"
        cmd = [
            mysql_cmd,
            "-h",
            os.getenv("DB_HOST", "127.0.0.1"),
            "-P",
            os.getenv("DB_PORT", "3306"),
            "-u",
            os.getenv("DB_USER", "root"),
            f"-p{os.getenv('DB_PASSWORD', '')}",
            "-N",
            "-e",
            query,
            db_name,
        ]
        try:
            output = subprocess.check_output(cmd, text=True, stderr=subprocess.DEVNULL)
            for line in output.splitlines():
                name, province, lat, lng = line.split("\t")
                rows.append({"name": name, "province": province, "lat": lat, "lng": lng})
        except (OSError, subprocess.CalledProcessError, ValueError):
            rows = []

    if not rows:
        rows = [
            {
                "name": city["name"],
                "province": city["province"],
                "lat": city["lat"],
                "lng": city["lng"],
            }
            for city in load_json(CITIES_JSON)[:35]
        ]

    with SEED_INPUTS.open("w", encoding="utf-8", newline="") as f:
        writer = csv.DictWriter(f, fieldnames=["name", "province", "lat", "lng"])
        writer.writeheader()
        writer.writerows(rows)
    print(f"wrote {len(rows)} rows to {SEED_INPUTS.relative_to(REPO_ROOT)}")


def gen_data(_: argparse.Namespace) -> None:
    existing = load_json(CITIES_JSON)
    cities = build_catalog(existing[:12])
    write_json(CITIES_JSON, cities)
    with SEED_INPUTS.open("w", encoding="utf-8", newline="") as f:
        writer = csv.DictWriter(f, fieldnames=["name", "province", "lat", "lng"])
        writer.writeheader()
        for city in cities:
            writer.writerow(
                {
                    "name": city["name"],
                    "province": city["province"],
                    "lat": city["lat"],
                    "lng": city["lng"],
                }
            )
    print(f"wrote {len(cities)} cities to {CITIES_JSON.relative_to(REPO_ROOT)}")


def gen_images(args: argparse.Namespace) -> None:
    cities = load_json(CITIES_JSON)
    achievements = load_json(ACHIEVEMENTS_JSON)
    print(f"gen-images offline={args.offline} pause={args.pause}", flush=True)
    manifest = []
    for index, city in enumerate(cities, start=1):
        manifest.extend(ensure_assets(REPO_ROOT, [city], [], offline=args.offline))
        write_manifest(REPO_ROOT, manifest)
        if index % 2 == 0 or index == len(cities):
            sourced_count = sum(
                1 for item in manifest
                if item.get("source") in {"Wikimedia Commons", "Wikimedia project page image"}
            )
            print(f"processed {index}/{len(cities)} cities, sourced assets={sourced_count}", flush=True)
        if not args.offline:
            time.sleep(args.pause)

    # Characters and badges are generated locally by design; badges are not city imagery.
    manifest.extend(ensure_assets(REPO_ROOT, [], achievements, offline=True))
    write_manifest(REPO_ROOT, manifest)
    print(f"wrote {len(manifest)} assets and scripts/asset_manifest.json")


def validate(_: argparse.Namespace) -> None:
    cities = load_json(CITIES_JSON)
    achievements = load_json(ACHIEVEMENTS_JSON)
    assert_static_files_exist(REPO_ROOT, cities, achievements)
    result = subprocess.run(
        ["go", "test", "./internal/seed"],
        cwd=REPO_ROOT / "backend",
        text=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.STDOUT,
    )
    if result.returncode != 0:
        sys.stdout.write(result.stdout)
        raise SystemExit(result.returncode)
    print(result.stdout.strip())
    print(f"validated {len(cities)} cities and {len(achievements)} achievements")


def main() -> None:
    parser = argparse.ArgumentParser(description="Build and validate demo seed data.")
    sub = parser.add_subparsers(required=True)
    handlers = {
        "audit-db": audit_db,
        "gen-data": gen_data,
        "gen-images": gen_images,
        "validate": validate,
    }
    for name, handler in handlers.items():
        cmd = sub.add_parser(name)
        if name == "gen-images":
            cmd.add_argument("--offline", action="store_true", help="skip Wikimedia downloads and generate local illustrations")
            cmd.add_argument("--pause", type=float, default=0.4, help="seconds to pause between city batches")
        cmd.set_defaults(func=handler)
    args = parser.parse_args()
    args.func(args)


if __name__ == "__main__":
    main()
