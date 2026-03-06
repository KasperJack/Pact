import json
from pathlib import Path
from .models import Package, Version
from .exceptions import PackageNotFoundError,InvalidManifestError, MissingKeyError

BUCKET_PATH = Path.cwd() / "bucket"



def load_package(package_name: str) -> Package:


    prefix = package_name[:2]
    index_file_path = BUCKET_PATH / prefix / package_name / "index.json"

    if not index_file_path.is_file():
        raise PackageNotFoundError("Package not found")


    try:
        with open(index_file_path, "r", encoding="utf-8") as f:
            data = json.load(f)
    except json.JSONDecodeError:
        raise InvalidManifestError("Invalid manifest file")

    #required keys
    required_keys = ["name", "slug", "release_year","igdb_id","default","ass"]
    validate_keys(data, required_keys, package_name)



    #default_version = data["default"]




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







def validate_keys(data: dict, required_keys: list[str], package_name: str):
    for key in required_keys:
        if key not in data:
            raise MissingKeyError(key, package_name)