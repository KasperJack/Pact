
#from pathlib import Path
from .loader import load_config

from .checker import Checker



class Pipeline:
    
    def __init__(self, config_file_path: str):
        self.config = load_config(config_file_path)
 


    def run(self):
        pass

        #checker = Checker(pacakge_raw)

   




