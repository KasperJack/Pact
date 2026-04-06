from pathlib import Path
import tempfile
import tomllib
from .models import Entity






import tomli_w
from datetime import date
from collections import defaultdict
from typing import Dict, List, Set





class PackageIndexer:
    def __init__(self, entities: list[Path]):
        self.entities = entities
        self.version_index: Dict[str, List[str]] = defaultdict(list)
        self.source_index: Dict[str, List[str]] = defaultdict(list)
        self.dlc_index: Dict[str, List[str]] = defaultdict(list)
        
        # Two ways for dates:
        # 1. For fast ID -> Date lookups
        self.id_to_date_lookup: Dict[str, str] = {} 
        # 2. For the sorted timeline view
        self.date_index: List[tuple] = []  # (date, package_id)

    def build_all_indexes(self):
        """Build all indexes from package TOML files"""
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
            self.id_to_date_lookup[package_id] = released
            
            # 4. DLC index
            content = data.get("content", {})
            dlcs = content.get("dlcs", [])

            if dlcs:
                for dlc in dlcs:
                    self.dlc_index[dlc].append(package_id)

        # Sort date index: Newest date first
        self.date_index.sort(key=lambda x: x[0], reverse=True)

    def get_release_date(self, package_id: str) -> str:
        """Instant lookup for an ID"""
        return self.id_to_date_lookup.get(package_id, "Not Found")

    def save_indexes(self, output_dir: Path):
        """Save all indexes to separate TOML files with grouped dates"""
        output_dir.mkdir(parents=True, exist_ok=True)

        # --- 1. Save Version & Source ---
        with open(output_dir / "version_index.toml", "wb") as f:
            tomli_w.dump({"versions": dict(self.version_index)}, f)
        
        with open(output_dir / "source_index.toml", "wb") as f:
            tomli_w.dump({"sources": dict(self.source_index)}, f)

        # --- 2. Save Date Index (Grouped Way) ---
        # We group by date here so the TOML is clean
        grouped_dates = defaultdict(list)
        for date_str, pkg_id in self.date_index:
            grouped_dates[date_str].append(pkg_id)

        # Sort keys to ensure the TOML file follows the timeline
        sorted_keys = sorted(grouped_dates.keys(), reverse=True)
        
        date_entries = []
        for date in sorted_keys:
            date_entries.append({
                "date": date,
                "packages": sorted(grouped_dates[date]) # Alphabetical pkgs per date
            })

        with open(output_dir / "date_index.toml", "wb") as f:
            tomli_w.dump({"entries": date_entries}, f)

        # --- 3. Save ID Lookup (For the search engine) ---
        # This is a small "Master List" for your CLI to check dates quickly
        with open(output_dir / "id_lookup.toml", "wb") as f:
            tomli_w.dump(self.id_to_date_lookup, f)

        with open(output_dir / "dlc_index.toml", "wb") as f:
            tomli_w.dump(dict(self.dlc_index), f)













class BaseBuilder:
    def __init__(self, backet_path: Path | str):
        self.bucket_path = Path(backet_path)



    def build_package(self,package_name: str):

        package_path = self._find_package(package_name)

        entities = self.get_available_entities(package_path)

        indexer = PackageIndexer(entities)
        indexer.build_all_indexes()
        indexer.save_indexes(Path("./indexes"))


        """
        with tempfile.TemporaryDirectory() as tmp_dir:

            tmp_path = Path(tmp_dir)

            for entity_path in entities:
                if not entity_path.is_file():
                    raise Exception(f"missing entity file at {entity_path}")

                try:
                    with open(entity_path, "rb") as f:
                        data = tomllib.load(f)
                except tomllib.TOMLDecodeError:
                    raise Exception(f"invalid entity file at {entity_path}")

                Entity.model_validate(data)
        """







    def _find_package(self, package_name: str) -> Path:

        prefix = package_name[:2]
        package_path = self.bucket_path / prefix / package_name
        
        if package_path.is_dir():
            return package_path
        else:
            raise Exception(f"pacakge not found at {package_path}")



    def get_available_entities(self,package_path: Path) -> list[Path]:
        entities_path = package_path / "entities"

        return [d / "e.toml" for d in entities_path.iterdir() if d.is_dir()]
    





    def _scan_dir(self, target_path: Path, ignore: list[str] | None = None) -> list[str]:

        to_ignore = ignore or []
        
        return [
            d.name for d in target_path.iterdir() 
            if d.is_dir() and d.name not in to_ignore
        ]
    


    def delete_old_indexes(self, package_path: Path):
        indexes_folder = package_path / "indexes"


        if not indexes_folder.is_dir():
            return
        
