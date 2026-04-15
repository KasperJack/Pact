
from datetime import date
from typing import Literal, Any,Union
from pydantic import BaseModel,ConfigDict,model_validator


class BoolNamespace(BaseModel):
    model_config = ConfigDict(extra="forbid")

    type: Literal["bool"]
    description: str
    reserved_flags: list[str]


class EnumNamespace(BaseModel):
    model_config = ConfigDict(extra="forbid")
    
    type: Literal["enum"]
    description: str


# union type
Namespace = BoolNamespace | EnumNamespace


class NamespaceRegistry(BaseModel):
    namespaces: dict[str, Namespace]

    def get(self, name: str) -> Namespace | None:
        return self.namespaces.get(name)




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



class EntityMeta(BaseModel):
    model_config = ConfigDict(extra="forbid")

    source:   str
    released: date
    version:  str




class Entity(BaseModel):
    id:                 str             
    meta:               EntityMeta
    public_interfaces:  PublicInterfaces  = PublicInterfaces()
    private_interfaces: PrivateInterfaces = PrivateInterfaces()


# ─── Entity Loader ────────────────────────────────────────────────────────────






class BoolInterfaceBlock(BaseModel):
    model_config = ConfigDict(extra="forbid")

    flags: list[str]
    default: Union[str, bool, None] = None

    @model_validator(mode="after")
    def validate_default(self):
        flags = self.flags
        default = self.default

        if len(flags) == 0:
            raise ValueError("no flags defineddd")


        if default is None:
            return self

        if len(flags) == 1:
            if not isinstance(default, bool):
                raise ValueError("default must be a boolean when only one flag is defined")

        else:
            if not isinstance(default, str):
                raise ValueError("default must be a string when multiple flags are defined")
            if default not in flags:
                raise ValueError("default must be one of the flags")

        return self
    




class EnumInterfaceBlock(BaseModel):
    model_config = ConfigDict(extra="forbid")

    keys: list[str]
    default: Union[str, bool, None] = None

    @model_validator(mode="after")
    def validate_default(self):
        keys = self.keys
        default = self.default

        if len(keys) == 0:
            raise ValueError("no flags defineddd")


        if default is None:
            return self

        if len(keys) == 1:
            if not isinstance(default, bool):
                raise ValueError("default must be a boolean when only one flag is defined")

        else:
            if not isinstance(default, str):
                raise ValueError("default must be a string when multiple flags are defined")
            if default not in keys:
                raise ValueError("default must be one of the flags")

        return self