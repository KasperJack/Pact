from __future__ import annotations
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from pathlib import Path






class PackageNotFoundError(Exception):
    exit_code = 4
    def __init__(self, path: str):
        super().__init__(f"pact: error: pacakge not found")
        self.path = path













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



class NamespaceFileNotFound(BaseError):
    def __init__(self, path:str):
        super().__init__(
            message = "Namespace file not found",
            path=path,
  
        )

class ReleaseFileNotFound(BaseError):
    def __init__(self, path:str):
        super().__init__(
            message = "release file not found",
            path=path,
  
        )


class ConfigParseError(BaseError):
    pass

class ConfigConversionError(BaseError):
    pass



class ErrorGroup(Exception):
    def __init__(self, errors: list[Exception]):
        self.errors = errors

    def __str__(self):
        lines = []
        for i, e in enumerate(self.errors, 1):
            lines.append(f"[{i}]-{e}")
        return "\n" + "\n".join(lines)