
from pathlib import Path
from .errors import PackageNotFoundError
from .parser import load_namespace_file, load_releases

from .checker import Checker
from typing import Any
from enum import Enum


class Folder(Enum):
    RELEASES = "releases"


class File(Enum):
    PACKAGE = "package"



class PackagePaths:
    def __init__(self, package_path: Path):
        self.root = package_path

    # folders
    @property
    def releases_folder(self) -> Path:
        return self.root / Folder.RELEASES.value



    # files
    @property
    def package_file(self) -> Path:
        return self.root / f"{File.PACKAGE.value}.conf"





class Compiler:
    
    def __init__(self, bucket_path: Path, pacakge_name: str):
        self.bucket_path = bucket_path
        self.pacakge_name = pacakge_name
        self.pacakge_path = self._find_package()
        self.paths = PackagePaths(self.pacakge_path)



    def run(self):

        pacakge_raw = self.get_package_data()

        checker = Checker(pacakge_raw)
        

        #print(namespace_raw)
        
        #releases_raw = self.get_releases_data()

        #checker.check_entities(entities_raw)






    def _find_package(self) -> Path:

        prefix = self.pacakge_name[:2]

        package_path = self.bucket_path / prefix / self.pacakge_name
        
        if package_path.is_dir():
            return package_path
        else:
            raise PackageNotFoundError(str(package_path))




    def get_package_data(self) -> dict[str,Any]:
        
        name_space_file = self.paths.package_file

        return load_namespace_file(name_space_file)



    def get_releases_data(self) -> dict[Path, dict]:

        releases_folder = self.paths.releases_folder
        paths =  [f / "release.conf" for f in releases_folder.iterdir() if f.is_dir()]

        return load_releases(paths)















"""
checker = Checker("libtree/namespaces.toml")

entities = checker.check([
    "libtree/release_1.6.toml",
    "libtree/release_1.7.toml",
    "libtree/release_1.8.toml",
])

# hand off to your indexer
indexer = Indexer(entities)
indexer.build()



def pipeline_run(bucket_path: Path ,pacakge_name: str) -> str:

    pacakge_path = _find_package(bucket_path,pacakge_name)

    # check if pacakge is already build 

    return "ass"
"""


