class InvalidManifestError(Exception):
    pass


class PackageNotFoundError(Exception):
    pass


class MissingKeyError(Exception):
    def __init__(self, key: str, package: str):
        super().__init__(f"Missing key '{key}' in package '{package}' manifest")
        self.key = key
        self.package = package