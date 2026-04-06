from pathlib import Path
import tempfile

class BaseBuilder:
    def __init__(self, backet_path: Path | str):
        self.bucket_path = Path(backet_path)



    def build_package(self,package_name: str):

        package_path = self._find_package(package_name)

        entities = self.get_available_entities(package_path)


        with tempfile.TemporaryDirectory() as tmp_dir:
            tmp_path = Path(tmp_dir)
            for e in entities:
                pass








    def _find_package(self, package_name: str) -> Path:

        prefix = package_name[:2]
        package_path = self.bucket_path / prefix / package_name
        
        if package_path.is_dir():
            return package_path
        else:
            raise Exception(f"pacakge not found at {self.bucket_path}")



    def get_available_entities(self,package_path: Path) -> list[Path]:
        entities_path = package_path / "entities"

        return [d / "entity.toml" for d in entities_path.iterdir() if d.is_dir()]
    





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
        
