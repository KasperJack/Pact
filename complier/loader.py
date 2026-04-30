from __future__ import annotations
from typing import Any

from pathlib import Path
from .errors import LoadError

from pact_core.loader import load
from pact_core.loader.exceptions import ParseError,ConfigConversionError



def load_config(path:str) -> dict[str,Any]:

    full_path = Path(path).resolve()

    if not full_path.exists():
        raise LoadError(f"Error: file not found: {full_path}") # L/E


    if not full_path.is_file():
        raise LoadError(f"Error: file not found: {full_path}")# L/E



    try:
        l = load(full_path)
    except ValueError as e:
        raise LoadError(str(e))
    except ParseError as e:
        raise LoadError(str(e))
    except ConfigConversionError as e:
        raise LoadError(str(e))
    return l