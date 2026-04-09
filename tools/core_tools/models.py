from __future__ import annotations
import tomllib
from pathlib import Path
from datetime import date
from typing import Literal, Any
from pydantic import BaseModel


# ─── Global Namespace Registry ────────────────────────────────────────────────

class Namespace(BaseModel):
    type: Literal["bool", "string", "enum"]
    description: str
    reserved_flags: list[str]


class NamespaceRegistry(BaseModel):
    namespaces: dict[str, Namespace]

    def get(self, name: str) -> Namespace | None:
        return self.namespaces.get(name)


def load_namespace_registry(path: str | Path) -> NamespaceRegistry:
    with Path(path).open("rb") as f:
        raw = tomllib.load(f)
    return NamespaceRegistry(namespaces={
        name: Namespace(**data)
        for name, data in raw.get("namespace", {}).items()
    })


# ─── Bool Interface ───────────────────────────────────────────────────────────

class BoolInterface(BaseModel):
    namespace: str              # only set on public interfaces
    local_var: str
    flags: list[str]
    default: str | bool | None
    description: str = ""       # inherited from namespace on public

    @property
    def full_path(self) -> str:
        return f"interface.{self.namespace}.{self.local_var}"


# ─── Interface Containers ─────────────────────────────────────────────────────

class PublicInterfaces(BaseModel):
    bool: dict[str, BoolInterface] = {}
    # string: dict[str, StringInterface] = {}   <- coming later
    # enum:   dict[str, EnumInterface]   = {}   <- coming later


class PrivateInterfaces(BaseModel):
    bool: dict[str, BoolInterface] = {}
    # string: dict[str, StringInterface] = {}   <- coming later
    # enum:   dict[str, EnumInterface]   = {}   <- coming later


# ─── Entity ───────────────────────────────────────────────────────────────────

class Entity(BaseModel):
    id: str
    source: str
    released: date
    version: str
    public_interfaces:  PublicInterfaces  = PublicInterfaces()
    private_interfaces: PrivateInterfaces = PrivateInterfaces()


# ─── Entity Loader ────────────────────────────────────────────────────────────

def load_entity(path: str | Path, registry: NamespaceRegistry) -> Entity:
    path = Path(path)

    if not path.exists():
        raise FileNotFoundError(f"Entity file not found: {path}")

    with path.open("rb") as f:
        raw = tomllib.load(f)  # tomllib raises TOMLDecodeError on duplicate keys

    public_bool:  dict[str, BoolInterface] = {}
    private_bool: dict[str, BoolInterface] = {}

    for key, group in raw.get("interface", {}).items():

        # ── PUBLIC: [interface.<namespace>.<local_var>] ──────────────────────
        # group is a dict whose values are also dicts (nested one more level)
        if isinstance(group, dict) and any(isinstance(v, dict) for v in group.values()):
            namespace_name = key
            ns = registry.get(namespace_name)

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





registry = load_namespace_registry("./type_check_test/namespaces.toml")
entity   = load_entity("./type_check_test/entity.toml", registry)

# Access a public bool interface
proto = entity.public_interfaces.bool.get("protocol")
if proto:
    print(proto.full_path)   # interface.download.protocol
    print(proto.flags)       # ['ssh', 'https']
    print(proto.default)     # 'https'
    print(proto.description) # inherited from namespace