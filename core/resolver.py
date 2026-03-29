#from pathlib import Path

from __future__ import annotations
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from .loader import Loader
    from pathlib import Path


from typing import Any # duck typing
from .exceptions import UserInputError



class resolver:
    def __init__(self, loader: Loader, package_path: Path, package__manifest: dict[str, Any], source: str, version: str, method: str):
        self.loader = loader # loader instance 
        self.package_path = package_path
        self.package__manifest = package__manifest

        # args can be None 
        self.target_source = source
        self.target_version = version
        self.target_method = method

        #list[str]
        #self.available_sources = []
        #self.available_versions = []
        #self.available_methods = []

        self.resolve()



    def resolve(self):

        #print(self.loader.package_name)

        available_sources = self.loader.get_available_sources(self.package_path)  

        if self.target_source: # is not None
            if self.target_source not in available_sources:
                raise UserInputError("source", self.target_source, available_sources)
        else:
            self.auto_select_source(available_sources)



        available_versions = self.loader.get_available_versions(self.package_path,self.target_source)

        if self.target_version:
            if self.target_version not in available_versions:
                raise UserInputError("version",self.target_version, available_versions)

        else:
            self.auto_select_version(available_versions)





    def auto_select_source(self, available_sources: list[str]):
        pref_sources = ["fitgirl","steamrip","steamunlocked"]

        if len(available_sources) == 1:
            self.target_source = available_sources[0]
            return

        default_source = self.package__manifest.get("preferred_s")

        if default_source:
            if default_source in available_sources:
                self.target_source = default_source
                return
            else:
                pass 
                #maybe log or warn

        for s in pref_sources:
            if s in available_sources:
                self.target_source = s
                return

        self.target_source = available_sources[0]

       
    

    def auto_select_version(self, available_versions: list[str]):

        if len(available_versions) == 1:
            self.target_version = available_versions[0]
            return

        default_version = self.package__manifest.get("default_version")

        if default:
            if "/" in default:
                version = default.split("/", 1)[1] 
                if version in available_versions:
                    self.target_version = version
                    return
                else:
                    pass 
                    #maybe log or warn
                    self.loader.package_name





        

       
    def get_available_methods(self):
        pass