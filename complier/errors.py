from __future__ import annotations
from typing import TYPE_CHECKING




from collections import defaultdict
import difflib


class LoadError(Exception):
    exit_code = 11
    def __init__(self, message: str):
        super().__init__(message)





class ConfigValidationError(Exception):
    exit_code = 21

    def __init__(self, msg: str):
        super().__init__(msg)

    @classmethod
    def from_pydantic_errors(cls, errors: list, filename: str = "config"):
        grouped = defaultdict(list)

        for err in errors:
            loc = list(err["loc"])

            # group (release, hash) → group=release, field=hash
            if len(loc) == 1:
                group = "root"
                field = loc[0]
            else:
                group = loc[0]
                field = loc[-1]

            msg = cls._format_error(err, field, group)
            grouped[group].append(msg)

        # build final output
        total = sum(len(v) for v in grouped.values())
        lines = [f"Found {total} errors in {filename}:\n"]

        for group in sorted(grouped.keys()):
            lines.append(f"  {group}:")
            for msg in grouped[group]:
                lines.append(f"    - {msg}")
            lines.append("")  # s

        return cls("\n".join(lines).strip())

    @staticmethod
    def _format_error(err, field, group):
        t = err["type"]

        if t == "missing":
            return f"missing field: {field}"

        elif t == "extra_forbidden":
            suggestion = ConfigValidationError._suggest_field(group, field)
            if suggestion:
                return f'unknown field: {field} (did you mean "{suggestion}"?)'
            return f"unknown field: {field}"

        elif t == "string_pattern_mismatch":
            return f"invalid field: {field} (invalid format)"

        else:
            return f"{field}: {err['msg']}"

    @staticmethod
    def _suggest_field(group, field):
        KNOWN_FIELDS = {
            "release": ["package_identifier", "url", "version", "hash"],
            "package": ["name", "package_identifier", "versioning", "license", "homepage", "description"],
        }

        options = KNOWN_FIELDS.get(group, [])
        matches = difflib.get_close_matches(field, options, n=1)

        return matches[0] if matches else None








class ConfigNotFoundError(LoadError):
    exit_code = 11
    def __init__(self, message: str):
        super().__init__(message)
  



















class BaseError(Exception):
    def __init__(self, message: str, *, cause: Exception | None = None, **context):
        self.message = message
        self.context = context
        self.cause = cause
        super().__init__(message)

    @property
    def error_type(self):
        return self.__class__.__name__

    def __str__(self):
        lines = [f"Error: {self.error_type}", " |"]

        for k, v in self.context.items():
            lines.append(f" | {k}: {v}")

        if self.cause:
            lines.append(f" | cause: {repr(self.cause)}")

        lines.append(" |")
        lines.append(f" | {self.message}")

        return "\n".join(lines)



