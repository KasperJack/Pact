import tomllib
from pathlib import Path
from .models import Package, Version
from.resolver import resolver
from .exceptions import MissingManifestError,InvalidManifestError,MissingKeyError
from dataclasses import dataclass
from typing import Any

BUCKET_PATH = Path.cwd() / "bucket -game-based"


class Loader:
    """
    All filesystem work lives here.
    Scans directories, reads TOML files, throws file-related errors.
    Uses Resolver for all decision logic.
    """
 
    def __init__(self, bucket_path: str | Path):
        self.bucket_path = Path(bucket_path)

    def load(
        self,
        package_name: str,
        source: str | None = None,
        version: str | None = None,
        method: str | None = None,
    ) -> Package:


        package_path = self._find_package(package_name)

        package__index_data = self._get_package__index(package_path, package_name)

        r = resolver(self, package_path,package__index_data,source,version,method)
        
        # call resover here ??
        target: InstallTarget = resolve(package_path,index_data,source,version,method)
        

        ## later fix all this part 
        if version is None:
            version = index_data["default_version"]
        else:
            pass


        

        version_file_path = package_path / version / "manifest.json"


        if not version_file_path.is_file():
            raise VersionManifestNotFoundError(package_name)


        try:
            with open(version_file_path, "r", encoding="utf-8") as f:
                version_data = json.load(f)
        except json.JSONDecodeError:
            raise InvalidVersionManifestError(package_name)


        validate_keys_version(version_data,package_name)


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




    def _find_package(self, package_name: str) -> str:

        prefix = package_name[:2]
        package_path = BUCKET_PATH / prefix / package_name
        
        if package_path.is_dir():
            return package_path
        else:
            pass
            #Raise pacakge not found 

    def _get_package__index(self, package_path: str, package_name: str) -> dict[str, Any]:
        index_file_path = package_path / "index.toml"

        if not index_file_path.is_file():
            raise IndexManifestNotFoundError(package_name)


        try:
            with open(index_file_path, "rb") as f:  #bytes
                index_data = tomllib.load(f)
        except tomllib.TOMLDecodeError:
            raise InvalidIndexManifestError(package_name)

        #required keys
        validate_keys_index(index_data,package_name)




    def validate_keys_index(data: dict, package_name: str):
        required_keys = ["name", "default_version"]
        for key in required_keys:
            if key not in data:
                raise MissingIndexKeyError(key, package_name)
        
        ids = data.get("ids", {})
        if "igdb" not in ids:
            raise MissingIndexKeyError("igdb", package_name)
            





    def validate_keys_version(data: dict, package_name: str):
        required_keys = ["name", "slug", "release_year","igdb_id","default","ass"]
        for key in required_keys:
            if key not in data:
                raise MissingVersionKeyError(key, package_name)
            




    def resolve(package_path,index_data, source, version, method) -> InstallTarget:
        # SOURCE
        sources = get_sources(package_path)

        if source:
            if source not in sources:
                raise SourceNotFound(source, sources)
        else:
            source = auto_select_source(index_data, sources)

        # VERSION
        versions = get_versions(package_path, source)

        if version:
            if version not in versions:
                raise VersionNotFound(version, versions)
        else:
            version = auto_select_version(index_data, source, versions)

        # METHOD
        methods = get_methods(package_path, source, version)

        if method:
            if method not in methods:
                raise MethodNotFound(method, methods)
        else:
            method = select_method(methods)


    ## error no avalable sources for the package
    def get_sources(package_path: str) -> list[str]:
        package_path = Path(package_path)

        return [
            p.name
            for p in package_path.iterdir()
            if p.is_dir() and p.name != "steam_builds"
        ]



    def auto_select_source(index_data: dict[str, Any], sources: list[str]) -> str:
        pref_sources = ["a"]
        default = index_data.get("default_version")

        if default:
            source, _ = default.split("/", 1) #!? 
            return source
        
        for s in pref_sources:
            if s in sources:
                return s

        return sources[0]






    ## error no avalable versions for the package
    def get_versions(package_path, source)-> list[str]:
        source_path = package_path / source
        return [
            p.name
            for p in source_path.iterdir()
            if p.is_dir()
        ]


