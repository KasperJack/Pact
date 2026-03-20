
""""
class RegistryError(Exception):
    pass


class ManifestError(RegistryError):
    pass





#missing manifest errors

class IndexManifestNotFoundError(ManifestError):
    def __init__(self, package: str):
        super().__init__(f"Index manifest not found for package '{package}'")
        self.package = package


class VersionManifestNotFoundError(ManifestError):
    def __init__(self, package: str):
        super().__init__(f"Version manifest not found for package '{package}'")
        self.package = package

#invalid JSON errors

class InvalidIndexManifestError(ManifestError):
    def __init__(self, package: str):
        super().__init__(f"Invalid index manifest for package '{package}'")
        self.package = package


class InvalidVersionManifestError(ManifestError):
    def __init__(self, package: str):
        super().__init__(f"Invalid version manifest for package '{package}'")
        self.package = package

#missing key errors

class MissingKeyError(ManifestError):
    def __init__(self, key: str, package: str):
        super().__init__(f"Missing key '{key}' in package '{package}' manifest")
        self.key = key
        self.package = package

class MissingIndexKeyError(MissingKeyError):
    pass


class MissingVersionKeyError(MissingKeyError):
    pass





"""


class PackageManagerError(Exception): 
    pass

class LoaderError(PackageManagerError): 
    pass
class InstallationError(PackageManagerError): 
    pass

class FileSystemError(PackageManagerError): 
    pass
class ValidationError(PackageManagerError):
    pass

class ResolutionError(LoaderError):
    """Base for when we can't figure out which files to use"""
    exit_code = 7

#loader
class MissingManifestError(LoaderError, FileSystemError):
    exit_code = 4
    def __init__(self, package_name: str, path: str, manifest_type: str = "package"):
        super().__init__(f"manifest not found for package '{package_name}'")
        self.package = package_name
        self.path = path
        self.manifest_type = manifest_type
    
class InvalidManifestError(LoaderError, FileSystemError):
    ## in the case of malformed json
    def __init__(self, package_name: str, path: str, manifest_type: str = "package"):
        super().__init__(f"Invalid manifest for package '{package_name}'")
        self.package_name = package_name
        self.path = path
        self.manifest_type = manifest_type

class MissingKeyError(LoaderError, ValidationError):
    exit_code = 5
    def __init__(self, key: str, package_name: str, path: str, manifest_type: str = "package"):
        super().__init__(f"Missing key '{key}' in package '{package_name}' manifest")
        self.key = key
        self.package_name = package_name
        self.path = path
        self.manifest_type = manifest_type


class UserInputError(ResolutionError):
    """The user asked for something specific that doesn't exist"""
    def __init__(self, item_type: str, requested: str, available: list):
        # item_type: "source", "version", or "method"
        super().__init__(f"{item_type.capitalize()} '{requested}' not found. Available: {', '.join(available)}")
        self.requested = requested
        self.available = available

class PackageEmptyError(ResolutionError):
    """The user didn't specify, but the system couldn't find a default"""
    def __init__(self, item_type: str, package_name: str):
        super().__init__(f"Could not auto-select {item_type} for '{package_name}'. No valid options found.")
        self.item_type = item_type