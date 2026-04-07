import tomllib
from pathlib import Path
from typing import Dict, List


class WildSearch:
    def __init__(self, package_path: Path):
        self.index_dir = package_path / "indexes"
        
        # Specialized Index Maps
        self.id_lookup: Dict[str, str] = {}         # ID -> Date
        self.version_map: Dict[str, List[str]] = {} # Version -> [IDs]
        self.source_map: Dict[str, List[str]] = {}  # Source -> [IDs]
        self.timeline: List[Dict] = []              # Sorted Release Ledger
        
        self._load_indexes()

    def _load_indexes(self):
        """Load all compiled indexes into memory"""
        # Load ID -> Date (The 'Master' list)
        self.id_lookup = self._load_toml("id_lookup.toml")
        
        # Load Versions (Grouped)
        v_data = self._load_toml("version_index.toml")
        self.version_map = v_data.get("versions", {})

        # Load Sources (Grouped)
        s_data = self._load_toml("source_index.toml")
        self.source_map = s_data.get("sources", {})

        # Load the Grouped Timeline
        t_data = self._load_toml("date_index.toml")
        self.timeline = t_data.get("entries", [])

    def _load_toml(self, filename: str) -> Dict:
        path = self.index_dir / filename
        if path.exists():
            with open(path, "rb") as f:
                return tomllib.load(f)
        return {}

    # --- Search Methods ---

    def find_by_version(self, version: str) -> List[str]:
        """Get all packages matching a specific version string"""
        return self.version_map.get(version, [])

    def find_by_source(self, source_name: str) -> List[str]:
        """Get all packages from a specific source (e.g., 'github', 'internal')"""
        # Normalize to lowercase if your indexer does the same
        return self.source_map.get(source_name.lower(), [])

    def get_package_info(self, package_id: str) -> Dict:
        """The 'Inspector': Pulls data from all indexes for one ID"""
        if package_id not in self.id_lookup:
            return {}
            
        return {
            "id": package_id,
            "release_date": self.id_lookup[package_id],
            "version": next((v for v, ids in self.version_map.items() if package_id in ids), "N/A"),
            "source": next((s for s, ids in self.source_map.items() if package_id in ids), "N/A")
        }
    



class Resolver:
    def __init__(
        self,
        package_path: Path,

        source: str | None = None,
        version: str | None = None,
        args: list[str]| None = None
    ):

        self.package_path = package_path
        

        self.source = source
        self.version = version


    def resolve(self) -> list[str]:

        w = WildSearch(self.package_path)
        if self.version:
            print(w.find_by_version(self.version))

        return ["f"]
    




    def _find_package(self) -> Path:

        prefix = self.package_name[:2]
        package_path = self.bucket_path / prefix / self.package_name
        
        if package_path.is_dir():
            return package_path
        else:
            raise Exception("Package not found")