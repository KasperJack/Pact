from core_tools import BaseBuilder











#bucket_path = "C:\\Users\\Aya\\Desktop\\game-get\\bucket\\bucket-entity-based"
#bb = BaseBuilder(bucket_path)
#bb.build_package("cyberpunk")


















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