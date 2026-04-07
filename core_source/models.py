from pydantic import BaseModel, field_validator #ValidationError
from datetime import date


class IdsModel(BaseModel):
    igdb: int          
    steam: int | None = None
    gog: int | None = None

class PackageManifest(BaseModel):
    name: str           
    release_year: int  
    preferred_source: str | None = None
    ids: IdsModel





######################################################
class VersionEntry(BaseModel):
    released: date


class RegistryManifest(BaseModel):
    preferred: str | None = None
    versions: dict[str, VersionEntry]

    @field_validator("versions")
    @classmethod
    def versions_not_empty(cls, v: dict) -> dict:
        if not v:
            raise ValueError("versions cannot be empty")
        return v







class ReleaseManifest(BaseModel):
    pass

class DownloadManifest(BaseModel):
    pass