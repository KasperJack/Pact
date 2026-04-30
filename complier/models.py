
from datetime import date
from typing import Literal, Any,Union,Annotated,List
from pydantic import BaseModel,ConfigDict,model_validator,field_validator,Field




class ConfigDef(BaseModel):
    model_config = ConfigDict(extra="forbid")  

    package: NewPackageDef | None = None
    release: NewReleaseDef






class NewReleaseDef(BaseModel):
    model_config = ConfigDict(extra="forbid")

    package_identifier: str | None = None
    url : str
    version: str
    hash: str


class NewPackageDef(BaseModel):
    model_config = ConfigDict(extra="forbid")

    name: str 
    package_identifier: str
    versioning: str

    license: str | None = None
    homepage: str | None = None
    description: str | None = None

