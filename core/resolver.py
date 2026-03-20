from pathlib import Path
from typing import Any
from .exceptions import UserInputError, PackageEmptyError



class resolver:
    def __init__(self, loader: Any, package_path: str, index_data: dict[str, Any], source: str, version: str, method: str):
        self.loader = loader # loader instance 
        self.package_path = package_path
        self.index_data = index_data

        # User input (can be None)
        self.target_source = source
        self.target_version = version
        self.target_method = method

        #list[str]
        #self.available_sources = []
        #self.available_versions = []
        #self.available_methods = []

        self.resolve()



    def resolve(self):

        self.loader.get_available_sources() ##throws an error no avalbale sources  

        if self.target_source:
            if self.target_source not in self.available_sources:
                raise UserInputError("source", self.target_source, self.available_sources)
        else:
            self.auto_select_source()


        self.loader.get_available_versions()

        if self.target_version:
            if self.target_version not in self.available_versions:
                raise UserInputError("version",self.target_version,self.available_versions)

        else:
            self.auto_select_version()








    def get_available_sources(self):

        package_path = Path(self.package_path)
        self.available_sources = [
        p.name
        for p in package_path.iterdir()
        if p.is_dir() and p.name != "steam_builds"
        ]

        if len(self.available_sources) == 0:
            raise PackageEmptyError("source",self.package_name)



    def auto_select_source(self):
        pref_sources = ["a","b","c"]

        if len(self.available_sources) == 1:
            self.target_source = self.available_sources[0]
            return

        default = self.index_data.get("default_version")

        if default:
            if "/" in default:
                source = default.split("/", 1)[0] 
                if source in self.available_sources:
                    self.target_source = source
                    return
                else:
                    pass 
                    #maybe log or warn


        for s in pref_sources:
            if s in self.available_sources:
                self.target_source = s
                return

        self.target_source = self.available_sources[0]


   





    def get_available_versions(self):
        versions_path = self.package_path / self.target_source

        self.available_versions = [
        d.name for d in versions_path.iterdir() if d.is_dir()
    ]

        if len(self.available_versions) == 0:
            raise PackageEmptyError("version",self.package_name)
        
    
    def auto_select_version(self):

        if len(self.available_versions) == 1:
            self.target_version = self.available_versions[0]
            return

        default = self.index_data.get("default_version")

        

       
    def get_available_methods(self):
        pass