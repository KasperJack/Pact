from __future__ import annotations
from typing import TYPE_CHECKING, Any, cast

if TYPE_CHECKING:
    from pathlib import Path

from pydantic import ValidationError



from .models import (
    ConfigDef,


)







class Checker:
    def __init__(self, configdef_raw: dict[str, Any]):

        self.configdef = self._check_configdef(configdef_raw)

  
 
    def _check_configdef(self, configdef_raw: dict) -> ConfigDef:
       
        cd = ConfigDef.model_validate(configdef_raw)

        if cd.pacakge == None:
            print(" this should be a release only defntion")
            # check thtat key target_package exists 
            # check that the target package exists in the remote/local repo
            # check the version against what is defined in the existing pacakge def 
            # check that the version does not exist already 

        else:
            print(" this should be a package and an intial release defention")
            # check package does not already exsit 
            # check init release def does not contain a target pacakge 
            # check versioning supported 
            # check init release uses samme versioning defined by pacakge 



        return cd



    # remove
    def validate_root_namespace(self,namespace_raw: dict):
        if "namespace" not in namespace_raw:
            raise ValueError("Missing 'namespace' key")

        if not isinstance(namespace_raw["namespace"], dict):
            raise TypeError("'namespace' must be a dictionary")

    # remove
    def validate_namespace_entries(self,namespace: dict):
        for name, value in namespace.items():
            if not isinstance(value, dict):
                raise TypeError(f"{name} must be a dictionary")


    # remove
    def validate_global_flags(self,objects: dict[str, OSBlock]):
        seen = set()
        for name, osb in objects.items():
            for flag in osb.reserved_flags:
                if flag in seen:
                    raise ValueError(f"Flag '{flag}' is duplicated across objects (found in {name})")
                seen.add(flag)