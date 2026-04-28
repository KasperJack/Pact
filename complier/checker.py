from __future__ import annotations
from typing import TYPE_CHECKING, Any, cast

if TYPE_CHECKING:
    from pathlib import Path

from pydantic import ValidationError



from .models import (
    OSBlock,


)






class Checker:
    def __init__(self, namespace_raw: dict[str, Any]):

        self.namespace_registry = self._check_namespace(namespace_raw)

  
 
    def _check_namespace(self, namespace_raw: dict) -> dict[str,OSBlock]:

        self.validate_root_namespace(namespace_raw)
        namespace = namespace_raw["namespace"]

        self.validate_namespace_entries(namespace)
        parsed: dict[str, OSBlock] = {}


        for name, block in namespace.items():

            
            try:
                parsed[name] = OSBlock.model_validate(block)
            except ValidationError as e:
                raise ValueError(f"{name}: invalid OSBlock → {e}")



        self.validate_global_flags(parsed)
        return parsed




    def validate_root_namespace(self,namespace_raw: dict):
        if "namespace" not in namespace_raw:
            raise ValueError("Missing 'namespace' key")

        if not isinstance(namespace_raw["namespace"], dict):
            raise TypeError("'namespace' must be a dictionary")


    def validate_namespace_entries(self,namespace: dict):
        for name, value in namespace.items():
            if not isinstance(value, dict):
                raise TypeError(f"{name} must be a dictionary")



    def validate_global_flags(self,objects: dict[str, OSBlock]):
        seen = set()
        for name, osb in objects.items():
            for flag in osb.reserved_flags:
                if flag in seen:
                    raise ValueError(f"Flag '{flag}' is duplicated across objects (found in {name})")
                seen.add(flag)