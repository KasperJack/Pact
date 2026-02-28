from dataclasses import dataclass
from typing import List, Optional
import json
from pathlib import Path
import sys


BUCKET_PATH = Path.cwd() / "bucket"


@dataclass
class Download:
    type: str
    url: str
    preferred: Optional[bool] = False

@dataclass
class Version:
    id: str
    version: str
    source: str
    size_mb: float
    notes: str
    downloads: List[Download]

@dataclass
class Package:
    name: str
    release_year: int
    igdb_id: int
    default: str
    versions: List[Version]



def load_package(package_name: str) -> Package:
    file_path = BUCKET_PATH / f"{package_name}.json"

    if not file_path.is_file():
        ## maybe sugest pakges with close name
        print('package does not exist') 
        sys.exit(1)


    try:
        with open(file_path, "r", encoding="utf-8") as f:
            data = json.load(f)
    except json.JSONDecodeError:
        print("Invalid JSON file.")
        sys.exit(2)




    versions = []
    for v in data["versions"]:
        downloads = [Download(**d) for d in v["downloads"]]
        versions.append(
            Version(
                id=v["id"],
                version=v["version"],
                source=v["source"],
                size_mb=float(v["size_mb"]),
                notes=v["notes"],
                downloads=downloads
            )
        )

    return Package(
        name=data["name"],
        release_year=data["realse_year"],  
        igdb_id=data["igdb_id"],
        default=data["default"],
        versions=versions
    )
