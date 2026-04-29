
from datetime import date
from typing import Literal, Any,Union,Annotated,List
from pydantic import BaseModel,ConfigDict,model_validator,field_validator,Field




class ConfigDef(BaseModel):
    model_config = ConfigDict(extra="forbid")  

    pacakge: NewPacakgeDef | None = None
    release: NewReleaseDef






class NewReleaseDef(BaseModel):
    model_config = ConfigDict(extra="forbid")

    target_package: str | None = None
    url : str
    version: str
    hash: str


class NewPacakgeDef(BaseModel):
    model_config = ConfigDict(extra="forbid")

    name: str 
    slug: str
    versioning: str

    license: str | None = None
    homepage: str | None = None
    description: str | None = None

