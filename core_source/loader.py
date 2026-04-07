import tomllib
from pathlib import Path

from pydantic import BaseModel, ValidationError
from .models import PackageManifest, RegistryManifest, ReleaseManifest, DownloadManifest

from .resolver import resolver
from .exceptions import MissingManifestError,InvalidManifestError,MissingKeyError,PackageNotFoundError, PackageEmptyError
from typing import Any, Callable, cast
from enum import Enum





class ManifestType(Enum):
    PACKAGE = "game"
    REGISTRY = "registry"
    RELEASE = "release"
    DOWNLOAD = "download"

    @property
    def filename(self) -> str:
        return f"{self.value}.toml"


    @property
    def model(self) -> type[BaseModel]:
        match self:
            case ManifestType.PACKAGE:
                return PackageManifest
            case ManifestType.REGISTRY:
                return RegistryManifest
            case ManifestType.RELEASE:
                return ReleaseManifest
            case ManifestType.DOWNLOAD:
                return DownloadManifest






class PackageComponent(Enum):
    SOURCES = "source"
    VERSIONS = "version"
    METHODS = "method"




class BaseLoader:

    """
    All filesystem work lives here.
    Scans directories, reads TOML files, throws file-related errors.
    """

    def __init__(self, bucket_path: str | Path, package_name: str):
        self.bucket_path = Path(bucket_path)
        self.package_name = package_name
        self.package_path = self._find_package(package_name)




    def _find_package(self, package_name: str) -> Path:

        prefix = package_name[:2]
        package_path = self.bucket_path / prefix / package_name
        
        if package_path.is_dir():
            return package_path
        else:
            raise PackageNotFoundError




    def load_package_manifest(self) -> PackageManifest:

        package_file_path = self.package_path / ManifestType.PACKAGE.filename
        return cast(PackageManifest, self._load_manifest(package_file_path, ManifestType.PACKAGE))
    




    def load_registry_manifest(self, source: str) -> RegistryManifest:

        registry_file_path = self.package_path / source / ManifestType.REGISTRY.filename
        return cast(RegistryManifest, self._load_manifest(registry_file_path, ManifestType.REGISTRY))




    def load_release_manifest(self, source: str, version: str) -> ReleaseManifest:

        release_file_path = self.package_path / source / version / ManifestType.RELEASE.filename
        return cast(ReleaseManifest, self._load_manifest(release_file_path, ManifestType.RELEASE))

    


    def load_download_manifest(self, source: str, version: str, method: str) -> DownloadManifest:

        download_file_path = self.package_path / source / version / method / ManifestType.DOWNLOAD.filename
        return cast(DownloadManifest, self._load_manifest(download_file_path, ManifestType.DOWNLOAD))





    def get_available_sources(self) -> list[str]:

        ignore_dir = ["steam_builds"]
        available_sources = self._scan_dir(self.package_path, ignore_dir)

        if len(available_sources) == 0:
            raise PackageEmptyError(self.package_path, PackageComponent.SOURCES)
        
        return available_sources






    def get_available_versions(self, source: str) -> list[str]:
        versions_path = self.package_path / source

        available_versions = self._scan_dir(versions_path)

        if len(available_versions) == 0:
            raise PackageEmptyError(versions_path, PackageComponent.VERSIONS)
        
        return available_versions




    def get_available_methods(self, source: str, version: str) -> list[str]:
        meathods_path = self.package_path / source / version

        available_methods = self._scan_dir(meathods_path)

        if len(available_methods) == 0:
            raise PackageEmptyError(meathods_path, PackageComponent.METHODS)
        
        return available_methods






    def _load_manifest(
        self,
        file_path: Path,
        manifest_type: ManifestType,
        #validator: Callable[[dict[str, Any], Path], None],
    ) -> BaseModel:
        

        if not file_path.is_file():
            raise MissingManifestError(file_path, manifest_type)

        try:
            with open(file_path, "rb") as f:
                data = tomllib.load(f)
        except tomllib.TOMLDecodeError:
            raise InvalidManifestError(file_path, manifest_type)

        #validator(data, manifest_type)
        #self._validate_keys(data,manifest_type)

        try:
            return manifest_type.model.model_validate(data)
        
        except ValidationError:
            #raise InvalidManifestError(file_path, manifest_type)
            raise






    def _scan_dir(self, target_path: Path, ignore: list[str] | None = None) -> list[str]:

        # If ignore is None, we treat it as an empty list to keep the logic clean
        to_ignore = ignore or []
        
        return [
            d.name for d in target_path.iterdir() 
            if d.is_dir() and d.name not in to_ignore
        ]
    




class TargetLoader(BaseLoader):

    def __init__(
        self,
        bucket_path: str | Path,
        package_name: str,
        source: str | None = None,
        version: str | None = None,
        method: str | None = None,
    ):
        super().__init__(bucket_path,package_name)

        self.source = source
        self.version = version
        self.method = method

    def load(self):

        

        package__manifest = self.load_package_manifest() ##raises an error "package not found"


        r = resolver(self, package__manifest, self.source, self.version, self.method)

        
        #print(r.target_source,r.target_version,r.target_method)
     
        #print(package__manifest.ids.gog)

        #rm = self.load_registry_manifest(self.source)
        #print(rm.versions["2.0.0.2"])        

 
