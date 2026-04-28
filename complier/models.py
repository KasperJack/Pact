
from datetime import date
from typing import Literal, Any,Union,Annotated,List
from pydantic import BaseModel,ConfigDict,model_validator,field_validator,Field




class OSBlock(BaseModel):
    model_config = ConfigDict(extra="forbid")  

    kind: Literal["Option", "Selection"]
    description: str
    reserved_flags: Annotated[List[str], Field(min_length=2, max_length=10)]

    @field_validator("reserved_flags")
    @classmethod
    def check_unique(cls, v):
        if len(v) != len(set(v)):
            raise ValueError("reserved_flags must contain unique items")
        return v