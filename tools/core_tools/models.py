from pydantic import BaseModel
from datetime import date



class Entity(BaseModel):
    id: str
    source: str
    released: date
    version: str
    