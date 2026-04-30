from __future__ import annotations
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from pathlib import Path



class LoadError(Exception):
    pass


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



