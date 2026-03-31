
from __future__ import annotations
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from pathlib import Path
    from .loader import ManifestType, PackageComponent




def _infer_package_name(path: Path) -> str:
    parts = path.parts
    for i, part in enumerate(parts):
        if len(part) == 2 and i + 1 < len(parts):
            return parts[i + 1]



class PackageManagerError(Exception): 
    pass

class LoaderError(PackageManagerError): 
    pass


class ResolutionError(LoaderError):
    """Base for when we can't figure out which files to use"""
    exit_code = 7

#loader

class PackageNotFoundError(LoaderError):
    pass


class MissingManifestError(LoaderError):
    exit_code = 4
    def __init__(self, path: Path, manifest_type: ManifestType):
        package_name = _infer_package_name(path)
        super().__init__(f" {manifest_type.value} manifest not found for package '{package_name}'")
        self.package = package_name
        self.path = path
        self.manifest_type = manifest_type
  
       
    
class InvalidManifestError(LoaderError):
    ## in the case of malformed toml
    def __init__(self, path: Path, manifest_type: ManifestType):
        package_name = _infer_package_name(path)
        super().__init__(f"Invalid {manifest_type.value} manifest for package '{package_name}'")
        self.package_name = package_name
        self.path = path
        self.manifest_type = manifest_type

class MissingKeyError(LoaderError):
    exit_code = 5
    def __init__(self, key: str, path: Path, manifest_type: ManifestType):
        package_name = _infer_package_name(path)
        super().__init__(f"Missing key '{key}' in package '{package_name}' {manifest_type.value} manifest")
        self.key = key
        self.package_name = package_name
        self.path = path
        self.manifest_type = manifest_type

class PackageEmptyError(LoaderError):
    def __init__(self, path: Path, item_type: PackageComponent):
        package_name = _infer_package_name(path)
        super().__init__(f"Could not select {item_type.value} for '{package_name}'. No options found")
        self.package_name = package_name
        self.path = path
        self.item_type = item_type





class UserInputError(ResolutionError):
    """The user asked for something specific that doesn't exist"""
    def __init__(self, item_type: str, requested: str, available: list):
        # item_type: "source", "version", or "method"
        super().__init__(f"{item_type.capitalize()} '{requested}' not found. Available: {', '.join(available)}")
        self.requested = requested
        self.available = available

