from core_tools import BaseBuilder




import tomllib
from pathlib import Path
from typing import Dict, List, Optional

class WildSearch:
    def __init__(self, index_dir: Path):
        self.index_dir = index_dir
        
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








#bucket_path = "C:\\Users\\Aya\\Desktop\\game-get\\bucket\\bucket-entity-based"
#bb = BaseBuilder(bucket_path)
#bb.build_package("cyberpunk")

w = WildSearch(Path("C:\\Users\\Aya\\Desktop\\game-get\\tools\\indexes"))

print(w.find_by_version("2.3"))
print(w.find_by_source("official"))



















"""

class PackageIndexer:
    def __init__(self, entities: list[Path]):
        self.entities = entities
        self.version_index: Dict[str, List[str]] = defaultdict(list)
        self.source_index: Dict[str, List[str]] = defaultdict(list)
        self.date_index: List[tuple] = []  # (date, package_id)
        self.dlc_index: Dict[str, List[str]] = defaultdict(list)
    
    def build_all_indexes(self):
        
        for toml_file in self.entities:
            with open(toml_file, "rb") as f:
                data = tomllib.load(f)
            
            Entity.model_validate(data)
            package_id = data.get("id")
    
            
            # 1. Version index
            version = data.get("version")

            self.version_index[version].append(package_id)
            
            # 2. Source index
            source = data.get("source")
            self.source_index[source].append(package_id)
            
            
            # 3. Release date index
            released = data.get("released")
            self.date_index.append((released, package_id))
            
            # 4. DLC index
            content = data.get("content", {})
            dlcs = content.get("dlcs", [])
            
            if dlcs:
                for dlc in dlcs:
                    self.dlc_index[dlc].append(package_id)

        
        # Sort date index
        self.date_index.sort(key=lambda x: x[0], reverse=True)  # newest first
    
    def save_indexes(self, output_dir: Path):
        
        output_dir.mkdir(exist_ok=True)
        
        # Save version index
        with open(output_dir / "version_index.toml", "wb") as f:
            tomli_w.dump(dict(self.version_index), f)
        
        # Save source index
        with open(output_dir / "source_index.toml", "wb") as f:
            tomli_w.dump(dict(self.source_index), f)
        
        # Save date index (as list of tables)
        date_index_data = [{"date": d, "packages": [pkg]} for d, pkg in self.date_index]
        with open(output_dir / "date_index.toml", "wb") as f:
            tomli_w.dump({"entries": date_index_data}, f)
        
        # Save DLC index
        with open(output_dir / "dlc_index.toml", "wb") as f:
            tomli_w.dump(dict(self.dlc_index), f)
"""