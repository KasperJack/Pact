from __future__ import annotations
from typing import Any

from pathlib import Path
from .errors import NamespaceFileNotFound,ConfigConversionError,ConfigParseError,ReleaseFileNotFound,ErrorGroup
import pyhocon # type: ignore
import json




def load_namespace_file(path: Path) -> dict:
    try:
        ns = load_hocon_file(path)
    except FileNotFoundError:
        raise NamespaceFileNotFound(str(path))

    return ns





def load_releases(paths: list[Path]) -> dict[Path,dict]:
    result = {}
    errors = []

    for p in paths:

        try:
            result[p] = load_hocon_file(p)
        except FileNotFoundError:
            errors.append(ReleaseFileNotFound(str(p)))
            #raise ReleaseFileNotFound(str(p))
            #raise ExceptionGroup()

    if errors:
        if len(errors) ==1:
            raise errors[0]
        
        raise ErrorGroup(errors=errors)

    return result   





def load_hocon_file(path: Path) -> dict[str, Any]:

    if not path.exists():
        raise FileNotFoundError(f"{path} does not exist")

    try:
        config = pyhocon.ConfigFactory.parse_file(path)
    except Exception as e:
        raise ConfigParseError("Failed to parse HOCON file",file=path, cause=e) from e

    try:
        return json.loads(json.dumps(config))
    except Exception as e:
        raise ConfigConversionError("Failed to convert config to dict",file=path, cause=e) from e





def _infer_release(path: Path) -> str:
    return path.parent.name

def _infer_package_name(path: Path) -> str:
    parts = path.parts
    for i, part in enumerate(parts):
        if len(part) == 2 and i + 1 < len(parts):
            return parts[i + 1]
    return "could not find package name from path"