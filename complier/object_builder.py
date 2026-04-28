from __future__ import annotations
import tomllib
from pathlib import Path
from typing import Any
from .parser import load_many,load_one

from .models import (
    Entity,
    Namespace,
    PublicInterfaces,
    PrivateInterfaces,
    BoolInterface,
    NamespaceRegistry,
)




class Checker:
    def __init__(self, namespace_path: Path):

        self.namespace_registry = self._load_namespace_registry(namespace_path)

        self.entities: list[Entity] = []
 



    def _load_namespace_registry(self, path: Path) -> NamespaceRegistry:
        with Path(path).open("rb") as f:
            raw = tomllib.load(f)
        return NamespaceRegistry(namespaces={
            name: Namespace(**data)
            for name, data in raw.get("namespace", {}).items()
        })





    def check(self,files: list[Path]) -> list[Entity]:
        
        for file in files:




    def load_entity(self,raw: dict[str, Any]) -> Entity:


        public_bool:  dict[str, BoolInterface] = {}
        private_bool: dict[str, BoolInterface] = {}

        for key, group in raw.get("interface", {}).items():

            # ── PUBLIC: [interface.<namespace>.<local_var>] ──────────────────────
            # group is a dict whose values are also dicts (nested one more level)
            if isinstance(group, dict) and any(isinstance(v, dict) for v in group.values()):
                namespace_name = key
                ns = self.namespace_registry.get(namespace_name)

                if ns is None:
                    raise ValueError(f"Unknown namespace: '{namespace_name}'")

                if ns.type != "bool":
                    continue  # string/enum handled later

                if len(group) > 1:
                    raise ValueError(
                        f"Namespace '{namespace_name}' used {len(group)} times — "
                        f"only one interface per namespace is allowed. "
                        f"Found: {list(group.keys())}"
                    )

                local_var, data = next(iter(group.items()))  # safe single unpack


                if any(isinstance(v, dict) for v in data.values()):
                    raise ValueError(
                        f"[interface.{namespace_name}.{local_var}] is too deep — "
                        f"expected [interface.<namespace>.<local_var>] but got extra nesting"
                    )
                
                if not data:
                    print("expected data in here ")

                expected_keys = ["default","flags"]
                for k in data:
                    if k not in expected_keys:
                        raise ValueError

                bad_flags = set(data["flags"]) - set(ns.reserved_flags)
                if bad_flags:
                    raise ValueError(
                        f"[interface.{namespace_name}.{local_var}] uses "
                        f"undeclared flags: {bad_flags}"
                    )

                public_bool[local_var] = BoolInterface(
                    namespace=namespace_name,
                    local_var=local_var,
                    flags=data["flags"],
                    default=data.get("default"),
                    description=ns.description,
                )

            # ── PRIVATE: [interface.<var_name>] ──────────────────────────────────
            else:
                # key = var_name, group = the interface data directly
                # coming later
                pass

        return Entity(
            id=raw["id"],
            source=raw["source"],
            released=raw["released"],
            version=raw["version"],
            public_interfaces=PublicInterfaces(bool=public_bool),
            private_interfaces=PrivateInterfaces(bool=private_bool),
        )