import json
import sys
from pathlib import Path
from .models import Package, Version
from .exceptions import PackageNotFoundError,InvalidManifestError

BUCKET_PATH = Path.cwd() / "bucket"


def load_package(package_name: str) -> Package:
    file_path = BUCKET_PATH / f"{package_name}.json"

    if not file_path.is_file():
        raise PackageNotFoundError("Package not found")


    try:
        with open(file_path, "r", encoding="utf-8") as f:
            data = json.load(f)
    except json.JSONDecodeError:
        raise InvalidManifestError("Invalid manifest file")




    versions = {}
    for v in data["versions"]:
        versions[v["id"]] = Version(
            id=v["id"],
            version=v["version"],
            source=v["source"],
            size_mb=float(v["size_mb"]),
            notes=v["notes"],
            downloads=v["downloads"]
        )

    return Package(
        name=data["name"],
        release_year=data["release_year"],
        igdb_id=data["igdb_id"],
        default=data["default"],
        versions=versions
    )
