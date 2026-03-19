from pathlib import Path
from typing import Any
from .exceptions import NotFoundError, AutoSelectError



class resolver:
    def __init__(self, package_name: str, package_path: str, index_data: dict[str, Any], source: str, version: str, method: str):
        self.package_name = package_name
        self.package_path = package_path
        self.index_data = index_data
        
        # User input (can be None)
        self.target_source = source
        self.target_version = version
        self.target_method = method

        self.available_sources = []
        self.available_versions = []
        self.available_methods = []

        self.resolve()



    def resolve(self):

        self.get_available_sources() ##throws an error no avalbale sources  

        if self.target_source:
            if self.target_source not in self.available_sources:
                raise NotFoundError("source", self.target_source, self.available_sources)
        else:
            self.auto_select_source()








    def get_available_sources(self):

        package_path = Path(self.package_path)
        self.available_sources = [
        p.name
        for p in package_path.iterdir()
        if p.is_dir() and p.name != "steam_builds"
        ]

        if len(self.available_sources) == 0:
            raise AutoSelectError("source",self.package_name)


    def auto_select_source(self):
        pref_sources = ["a","b","c"]
        default = self.index_data.get("default_version")

        if default:
            source, _ = default.split("/", 1) #!? 
            return source
        
        for s in pref_sources:
            if s in sources:
                return s

   





    def get_avalable_versions(self):
        pass
    def get_avalable_methods(self):
        pass