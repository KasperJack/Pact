from dataclasses import dataclass


"""
@dataclass
class Download:
    type: str
    url: str
    
"""



@dataclass
class Version:
    id: str
    version: str
    source: str
    size_mb: float
    notes: str
    downloads: dict[str, list[str]]


@dataclass
class Package:
    name: str
    release_year: int
    igdb_id: int
    default: str
    versions: dict[str, Version]