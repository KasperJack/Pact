from pydantic import BaseModel #ValidationError



class IdsModel(BaseModel):
    igdb: int           # required
    steam: int | None = None
    gog: int | None = None

class PackageManifest(BaseModel):
    name: str           # required
    release_year: int   # required
    preferred_source: str | None = None
    ids: IdsModel

class RegistryManifest(BaseModel):
    pass

class ReleaseManifest(BaseModel):
    pass

class DownloadManifest(BaseModel):
    pass