
from pathlib import Path
from .errors import PackageNotFoundError
from .parser import load_entities,load_namespace_file

from .checker import Checker
from typing import Any
from enum import Enum


class Folder(Enum):
    ENTITIES = "entities"


class File(Enum):
    NAMESPACE = "namespace"



class PackagePaths:
    def __init__(self, package_path: Path):
        self.root = package_path

    # folders
    @property
    def entities_folder(self) -> Path:
        return self.root / Folder.ENTITIES.value



    # files
    @property
    def namespace_file(self) -> Path:
        return self.root / f"{File.NAMESPACE.value}.toml"





class Compiler:
    
    def __init__(self, bucket_path: Path, pacakge_name: str):
        self.bucket_path = bucket_path
        self.pacakge_name = pacakge_name
        self.pacakge_path = self._find_package()
        self.paths = PackagePaths(self.pacakge_path)



    def run(self):

        #namespace_raw = self.get_namespace_data()
        #print(namespace_raw)

        entities_raw = self.get_entities_data()

        #checker = Checker(namespace_raw)

        #checker.check_entities(entities_raw)






    def _find_package(self) -> Path:

        prefix = self.pacakge_name[:2]

        package_path = self.bucket_path / prefix / self.pacakge_name
        
        if package_path.is_dir():
            return package_path
        else:
            raise PackageNotFoundError(package_path)




    def get_namespace_data(self) -> dict[str,Any]:
        
        name_space_file = self.paths.namespace_file

        return load_namespace_file(name_space_file)



    def get_entities_data(self) -> dict[Path, dict]:

        entities_folder = self.paths.entities_folder
        paths =  [f / "entity.toml" for f in entities_folder.iterdir() if f.is_dir()]

        return load_entities(paths)















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


