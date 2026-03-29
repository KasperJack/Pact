import tomllib
from pathlib import Path
from .models import Package, Version
from.resolver import resolver
from .exceptions import MissingManifestError,InvalidManifestError,MissingKeyError,PackageNotFoundError, PackageEmptyError
from typing import Any


class Loader:
    """
    All filesystem work lives here.
    Scans directories, reads TOML files, throws file-related errors.
    Uses Resolver for all decision logic.
    """
 
    def __init__(self, bucket_path: str | Path):
        self.bucket_path = Path(bucket_path)
        self.package_name = None

    def load(
        self,
        package_name: str,
        source: str | None = None,
        version: str | None = None,
        method: str | None = None,
    ) -> Package:


        package_path = self._find_package(package_name) ##raise package not found error 

        package__manifest = self._load_package_manifest(package_path) ##raises an error "package not found"

        #print(package__index_data)

        r = resolver(self, package_path, package__manifest, source, version, method)

   
        #print(r.target_source,r.target_version,r.target_method)
        #return
    
        ## TODO: remove all this part later *
        ## figure out where to do the version , method valibation 
        ## build the data class that holds the resolved pacakge 
        ## figure out how to get metadata without building a full package instance 
        ##  add type checking for mainfest data 




    def _find_package(self, package_name: str) -> Path:

        prefix = package_name[:2]
        package_path = self.bucket_path / prefix / package_name
        
        if package_path.is_dir():
            self.package_name = package_name
            return package_path
        else:
            raise PackageNotFoundError


    def _load_package_manifest(self, package_path: Path) -> dict[str, Any]:
        package_file_path = package_path / "game.toml"

        if not package_file_path.is_file():
            raise MissingManifestError(self.package_name,package_path,"game")


        try:
            with open(package_file_path, "rb") as f:  #bytes
                manifest_data = tomllib.load(f)
        except tomllib.TOMLDecodeError:
            raise InvalidManifestError(self.package_name,package_file_path,"game")

        #required keys
        self.validate_keys_index(manifest_data,package_file_path)
        
        return manifest_data



    def get_available_sources(self, package_path: str) -> list[str]:

        package_path = Path(package_path)
        available_sources = [
        p.name
        for p in package_path.iterdir()
        if p.is_dir() and p.name != "steam_builds"
        ]

        if len(available_sources) == 0:
            raise PackageEmptyError(self.package_name,package_path,"source")
        
        return available_sources


    def get_available_versions(self, package_path: str, source: str) -> list[str]:
        versions_path = package_path / source

        available_versions = [
        d.name for d in versions_path.iterdir() if d.is_dir()
    ]

        if len(available_versions) == 0:
            raise PackageEmptyError(self.package_name,versions_path,"version")
        
        return available_versions



    def get_available_methods(self, package_path: str, source: str, version: str) -> list[str]:
        meathods_path = package_path / source / version

        available_methods = [
        d.name for d in meathods_path.iterdir() if d.is_dir()
    ]

        if len(available_methods) == 0:
            raise PackageEmptyError(self.package_name,meathods_path,"method")
        
        return available_methods




    def load_registry_manifest(self, package_path: str, source: str) -> dict[str, Any]:

        registry_file_path = package_path / source / "registry.toml"

        if not registry_file_path.is_file():
            raise MissingManifestError(self.package_name, registry_file_path, "registry")


        try:
            with open(registry_file_path, "rb") as f:  #bytes
                index_data = tomllib.load(f)
        except tomllib.TOMLDecodeError:
            raise InvalidManifestError(self.package_name,registry_file_path,"index")

        #required keys
        self.validate_keys_index(index_data,registry_file_path)
        
        return index_data





    def validate_keys_index(self,data: dict, index_file_path: str):
        required_keys = ["name", "default_version"]
        for key in required_keys:
            if key not in data:
                raise MissingKeyError(key,self.package_name,index_file_path,"index")
        
        ids = data.get("ids", {})
        if "igdb" not in ids:
            raise MissingKeyError("igdb",self.package_name,index_file_path,"index")
            





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





