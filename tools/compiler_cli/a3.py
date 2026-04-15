from __future__ import annotations
from typing import TYPE_CHECKING, Any

if TYPE_CHECKING:
    from pathlib import Path









from .models import (
    Entity,
    EntityMeta,
    Namespace,
    BoolNamespace,
    EnumNamespace,

    PublicInterfaces,
    PrivateInterfaces,
    BoolInterface,
    NamespaceRegistry,


    BoolInterfaceBlock,

)



def _build_namespace(name: str, data: dict) -> BoolNamespace | EnumNamespace:

    if any(isinstance(v, dict) for v in data.values()):
        raise Exception(
            f"Namespace '{name}' is too deep — "
            f"expected [namespace.<name>] not [namespace.<name>.<extra>]"
        )


    match data.get("type"):
        case "bool":
            return BoolNamespace(**data)
        case "enum":
            return EnumNamespace(**data)
        case None:
            raise Exception(f" namespace type not defined for {name} block'")
        case _:
            raise Exception(f"Unknown namespace type: '{data.get('type')}'")




class Checker:
    def __init__(self, namespace_raw: dict[str, Any]):

        self.namespace_registry = self._load_namespace_registry(namespace_raw)

        #self.entities: list[Entity] = []
 





    def _load_namespace_registry(self, namespace_raw: dict) -> NamespaceRegistry:
        allowed_keys = {"namespace"}
        unexpected   = set(namespace_raw.keys()) - allowed_keys

        if unexpected:
            raise Exception(f"Unexpected keys in namespace file: {unexpected}")

        return NamespaceRegistry(namespaces={
            name: _build_namespace(name,data)
            for name, data in namespace_raw.get("namespace", {}).items()
        })





    def check_entities(self, entities: dict[Path, dict])-> list[Entity]:
        result = []
        for path, raw in entities.items():
            entity = self._build_entity(path, raw)
            result.append(entity)
        return result





    def _build_entity(self,path:Path, raw: dict[str, Any]) -> Entity:

        allowed_keys = {"interface","meta"}
        unexpected   = set(raw.keys()) - allowed_keys

        if unexpected:
            raise Exception(f"Unexpected keys in entity file: {unexpected}")


      
        meta_raw = raw.get("meta")
        if meta_raw is None:
            raise Exception(path, "Missing [meta] block") #E:  InvalidEntitySchema


        #meta = EntityMeta(**meta_raw) ##E: expect  pydatic errors 
        meta = EntityMeta.model_validate(meta_raw) #E: expect  pydatic errors

        #print(path.parent.name)

        #folder_name = path.parent.name
        #if folder_name != meta.id:
        #    raise Exception(path, folder_name, meta.id) ##E: EntityIdMismatch






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
                
                #if not data:
                #    raise Exception(f"[interface.{namespace_name}.{local_var}] is empty — expected flags and default") #E: raise InvalidEntitySchema


                bi = BoolInterfaceBlock.model_validate(data)




                bad_flags = set(bi.flags) - set(ns.reserved_flags)
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

        raise NotImplementedError("bop")
        return Entity(
            id=path.parent.name,
            meta=meta,
            public_interfaces=PublicInterfaces(bool=public_bool),
            private_interfaces=PrivateInterfaces(bool=private_bool),
        )
    


    ## use this fucn insted 
    def check_public_bool_interface(self,namespace_name: str, ns: BoolNamespace, group: dict) -> BoolInterface: #-> tuple[str, BoolInterface]:
                
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
                
               

                bi = BoolInterfaceBlock.model_validate(data)




                bad_flags = set(bi.flags) - set(ns.reserved_flags)
                if bad_flags:
                    raise ValueError(
                        f"[interface.{namespace_name}.{local_var}] uses "
                        f"undeclared flags: {bad_flags}"
                    )
                

                return BoolInterface(
                    namespace=namespace_name,
                    local_var=local_var,
                    flags=bi.flags,
                    default=bi.default,
                    description=ns.description,
                )