
from pathlib import Path
from .parser import load_namespace_file, load_releases

from .checker import Checker
from pact_core.loader import load


class Pipeline:
    
    def __init__(self, config_file_path: Path):
        self.config_file_path = config_file_path
 


    def run(self):

        #pacakge_raw = self.get_package_data()

        #checker = Checker(pacakge_raw)
        l = load(self.config_file_path)
        #
        




















