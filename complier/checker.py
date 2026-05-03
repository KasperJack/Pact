from __future__ import annotations
from typing import TYPE_CHECKING, Any

from pathlib import Path

from .errors import ConfigValidationError

from pact_core.models import NewPackageConfig,NewVersionConfig,ValidationError
from pact_core.repo_client import LocalFSTransport,RepoClient











class Checker:
    def __init__(self, configdef_raw: dict[str, Any]):

        self.configdef = self._check_configdef(configdef_raw)

        #print("gg")
 




    def _check_configdef(self, configdef_raw: dict) -> NewPackageConfig | NewVersionConfig:


        if "package" in configdef_raw:

            self._validate_new_package(configdef_raw)
        else:
            self._validate_new_version(configdef_raw)


        rc = RepoClient(LocalFSTransport(Path("C:\\Users\\Aya\\Documents\\projects\\pact\\Pact-ci\\test-buckets\\defult\\cy\\cy2")))

        rc.test_check_file(Path("ts"))

      
            #print(" this should be a release only defntion")
            # check thtat key target_package exists 
            # check that the target package exists in the remote/local repo (new remote/local repo status moduel in core_lib )
            # check the version against what is defined in the existing pacakge def 
            # check that the version does not exist already 

  
            #print(" this should be a package and an intial release defention")
            # check package does not already exsit 
            # check init release def does not contain a target pacakge 
            # check versioning supported 
            # check init release uses samme versioning defined by pacakge 



        return 




    def _validate_new_package(self,config: dict):

        try:
       
            np_config = NewPackageConfig.model_validate(config)
        except ValidationError as e:

            #raise ConfigValidationError.from_exception_data(str(e))
            #raise ConfigValidationError.from_pydantic_errors(e.errors(include_url=False))
            raise ConfigValidationError.from_pydantic_errors(e.errors(include_url=False))





    def _validate_new_version(self,config:dict):


        try:
       
            np_config = NewVersionConfig.model_validate(config)
        except ValidationError as e:

            #raise ConfigValidationError.from_exception_data(str(e))
            #raise ConfigValidationError.from_pydantic_errors(e.errors(include_url=False))
            raise ConfigValidationError.from_pydantic_errors(e.errors(include_url=False))