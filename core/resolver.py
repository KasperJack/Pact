from pathlib import Path
from typing import Any




class resolver:
    def __init__(self, package_name: str, package_path: str, index_data: dict[str, Any], source: str, version: str, method: str):
        self.package_name = package_name
        self.package_path = package_path
        self.index_data = index_data
        self.target_source = source
        self.target_version = version
        self.target_method = method

        self.avalable_sources = []
        self.avalable_versions = []
        self.avalable_methods = []

        self.resolve()



    def resolve(self):

        self.get_avalable_sources() ##throws an error no avalbale sources  

        if self.target_source:
            if self.target_source not in self.avalable_sources:
                raise SomeError(self) #" fix this"
        else:
            self.auto_select_source()








    def get_avalable_sources(self):

        package_path = Path(self.package_path)
        self.avalable_sources = [
        p.name
        for p in package_path.iterdir()
        if p.is_dir() and p.name != "steam_builds"
    ]



    def get_avalable_versions(self):
        pass
    def get_avalable_methods(self):
        pass