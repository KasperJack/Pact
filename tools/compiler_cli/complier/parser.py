from __future__ import annotations
from typing import Any

import tomllib
from pathlib import Path
from .errors import NamespaceFileNotFound,EntityNotFound

def load_namespace_file(path: Path) -> dict:
    if not path.exists():
        raise NamespaceFileNotFound(path)
    with path.open("rb") as f:
        return tomllib.load(f)


def load_entities(paths: list[Path]) -> dict[Path, dict]:
    results = {}
    for p in paths:
        p = Path(p)
        if not p.exists():
            raise EntityNotFound(p)
        with p.open("rb") as f:
            results[p] = flatten_entity(tomllib.load(f))
    return results






















def flatten_entity(raw: dict[str, Any]) -> dict:
    result = {
        "options":    {},
        "selections": {},
    }


    seen_local_var: dict[str, str] = {}  # local_var -→ where it was first seen
    seen_parents: dict[str, str] = {}

    def register(local_var: str, source: str) -> None:
        if local_var in seen_local_var:
            raise ValueError(
                f"var '{local_var}' already declared under '{seen_local_var[local_var]}' "
                f"— cannot reuse in '{source}'"
            )
        seen_local_var[local_var] = source




    def register_parent(parent: str, source: str) -> None:
        if parent in seen_parents:
            raise ValueError(
                f"namespace '{parent}' already declared as an '{seen_parents[parent]}' "
                f"— cannot reuse as '{source}'"
            )
        seen_parents[parent] = source









    # ── options: [option.<parent>.<local_var>] ────────────────────────────
    for parent, group in raw.get("option", {}).items():

        if not isinstance(group, dict) or not any(isinstance(v, dict) for v in group.values()):
            raise ValueError(f"[option.{parent}] is malformed — expected [option.<parent>.<local_var>]")

        if len(group) > 1:
            raise ValueError(f"[option.{parent}] already declared — only one allowed")

        register_parent(parent, "option")

        local_var, data = next(iter(group.items()))

        if any(isinstance(v, dict) for v in data.values()):
            raise ValueError(f"[option.{parent}.{local_var}] is too deep")
        
        if local_var != "_":
            register(local_var,f"option.{parent}")

        if not data:
            raise ValueError(f"[option.{parent}.{local_var}] is empty")

        unexpected = set(data.keys()) - {"flags", "default"}
        if unexpected:
            raise ValueError(f"[option.{parent}.{local_var}] unexpected keys: {unexpected}")

        if "flags" not in data:
            raise ValueError(f"[option.{parent}.{local_var}] missing 'flags'")

        result["options"][parent] = {
            "local_var": local_var,
            "flags":     data["flags"],
            "default":   data.get("default"),
        }



    # ── selections: [selection.<parent>.<local_var>] ──────────────────────
    for parent, group in raw.get("selection", {}).items():

        if not isinstance(group, dict) or not any(isinstance(v, dict) for v in group.values()):
            raise ValueError(f"[selection.{parent}] is malformed — expected [selection.<parent>.<local_var>]")

        if len(group) > 1:
            raise ValueError(f"[option.{parent}] already declared — only one allowed")

        register_parent(parent, "selection")

        local_var, data = next(iter(group.items()))

        if any(isinstance(v, dict) for v in data.values()):
            raise ValueError(f"[selection.{parent}.{local_var}] is too deep")
        

        if local_var != "_":
            register(local_var,f"selection.{parent}")

        if not data:
            raise ValueError(f"[selection.{parent}.{local_var}] is empty")

        unexpected = set(data.keys()) - {"flags", "default"}
        if unexpected:
            raise ValueError(f"[selection.{parent}.{local_var}] unexpected keys: {unexpected}")

        if "flags" not in data:
            raise ValueError(f"[selection.{parent}.{local_var}] missing 'flags'")

        result["selections"][parent] = {
            "local_var": local_var,
            "flags":     data["flags"],
            "default":   data.get("default", []),
        }

    return result