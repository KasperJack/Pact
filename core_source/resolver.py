#from pathlib import Path

from __future__ import annotations
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from .loader import TargetLoader
    from .models import PackageManifest, RegistryManifest


from .exceptions import UserInputError



class resolver:
    def __init__(self, loader: TargetLoader, package__manifest: PackageManifest, source: str, version: str, method: str):
        self.loader = loader # loader instance 
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

        available_sources = self.loader.get_available_sources()  

        if self.target_source: # is not None
            if self.target_source not in available_sources:
                raise UserInputError("source", self.target_source, available_sources)
        else:
            self.auto_select_source(available_sources)



        available_versions = self.loader.get_available_versions(self.target_source)

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


        if self.package__manifest.preferred_source:
            if self.package__manifest.preferred_source in available_sources:
                self.target_source = self.package__manifest.preferred_source
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

        registry_manifest = self.loader.load_registry_manifest(self.target_source)

        
        if registry_manifest.preferred:
            if registry_manifest.preferred in available_versions:
                self.target_version = registry_manifest.preferred
                return
            else:
                pass

        latest_ver = max(
            registry_manifest.versions.items(),
            key=lambda item: item[1].released
        )[0]
        
        print(latest_ver)


        
