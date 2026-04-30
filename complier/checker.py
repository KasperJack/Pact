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

        if cd.package == None:
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



