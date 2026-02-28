from dataclasses import dataclass
from typing import List, Optional

@dataclass
class Download:
    type: str
    url: str
    preferred: Optional[bool] = False

@dataclass
class Version:
    id: str
    version: str
    source: str
    size_mb: float
    notes: str
    downloads: List[Download]

@dataclass
class Package:
    name: str
    release_year: int
    igdb_id: int
    default: str
    versions: List[Version]