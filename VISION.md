

Software on Windows is a mess. There is no standard model for what an "installed package" actually means. Package managers like Scoop, Chocolatey, and winget all treat packages as a delivery mechanism — they care about getting software onto your machine, not about what it looks like after.

The result: you never really know what you have installed, who owns it, or whether you can cleanly remove it.

Pact is an attempt to fix that.

---

## Core Idea: Software Has Types

Pact defines three models for how software exists on a Windows system. Every package must declare which one it is.

### portable

> Just files in a directory. Nothing else.

- No installer runs
- No ARP entry
- No registry writes
- No files outside its own directory
- No auto-update

Pact controls the version. Multiple versions can coexist. Uninstall means deleting the folder.

---

### managed

> Runs an installer. Pact tracks it.

- Runs an `.exe` or `.msi` installer
- Writes to ARP (Add/Remove Programs)
- May scatter files across `Program Files`, `AppData`, `System32`
- No auto-update

Pact controls the version. One version at a time — the installer replaces the previous one. Uninstall means running the uninstaller.

---

### self_managed

> Pact bootstraps it. The software takes over.

- Pact runs the installer once
- The software has its own background update mechanism
- Pact no longer controls the version after install
- No point in version pinning or upgrading via Pact

Uninstall means running their uninstaller — if it works.

---

### Why This Matters

Declaring a type gives Pact a model of the software. It knows what it owns, what it can guarantee, and what it can't. Most package managers don't make this distinction — everything is treated the same, which is why uninstalls break things and version pinning is unreliable.

> These types are not 100% enforced. Starlark hooks can extend or override behavior at any point — before/after install, uninstall, version switch, or migration. You can also write a package with no script at all if the software fits the type completely. The types are still being refined.

---

## Architecture

```
CLI (runner)
    └── provides Repository + LocalState interfaces
        └── Core (orchestration + Starlark runtime)
            └── Platform layer (filesystem, registry, shortcuts, junctions)
                └── System mutations (controlled, auditable)
```

Core is functionally pure — deterministic Starlark logic isolated from side effects through narrow interfaces. The platform layer is the only place system mutations happen, making every operation traceable.

---

## Package Definition

Packages are defined in two layers:

**Global definition** (`windirstat.hcl`) — applies to all versions:
```hcl
identifier  = "windirstat"
name        = "WinDirStat"
kind        = "portable"
versioning  = "semver"
description = "Disk usage viewer"
homepage    = "https://windirstat.net"
license     = "GPL-2.0"

shortcut {
    name = "WinDirStat"
    exe  = "WinDirStat.exe"
    icon = "WinDirStat.exe"
}
```

**Release definition** (`release.hcl`) — per version, has download URLs and hashes:
```hcl
version = "2.6.1"

source {
    x64   { url = "..." sha256 = "..." }
    x86   { url = "..." sha256 = "..." }
    arm64 { url = "..." sha256 = "..." }
}
```

**Install script** (`script.star`) — Starlark, optional, for anything non-standard:
```python
def install():
    extract(dist(), staging())
    if os.x64():
        move(path.join(staging(), "x64/*"), dir())
```

If no script is defined, Pact handles the install automatically based on the package type.

---

## Entry Points

Pact exposes software to the user through typed entry points declared in the global definition.

```hcl
# GUI app — creates a desktop shortcut
shortcut {
    name = "WinDirStat"       # optional, falls back to package name
    exe  = "WinDirStat.exe"
    icon = "WinDirStat.exe"   # optional, falls back to exe
    args = ""                  # optional
}

# CLI tool — creates a shim on PATH
command {
    name = "rg"
    exe  = "rg.exe"
}
```

- `shortcut` → creates a `.lnk` on the desktop pointing to the active version
- `command` → creates a shim on PATH that forwards to the active version

Shortcuts and shims always point to a `current/` junction. Switching versions updates the junction — entry points never change.

---

## Version Management

Each package has a `current/` junction that always points to the active version:

```
C:/pact/packages/windirstat/
├── current/          ← junction, points to active version
├── 2.6.1/
└── 2.7.0/
```

Switching versions = update one junction. Shortcuts and shims are never recreated.

---

## Starlark Hooks

Scripts can hook into any lifecycle event:

```python
def pre_install():  pass
def post_install(): pass
def pre_remove():   pass
def post_remove():  pass
def pre_switch():   pass
def post_switch():  pass
def on_update():    pass

# if install() is defined, Pact does nothing automatically
def install():      pass
```

---

## Roadmap

**Phase 1 — Foundation**
- Local filesystem repository backend
- Install pipeline for `portable` packages
- Shortcut and shim creation
- Lockfile tracking

**Phase 2 — Infrastructure**
- Git-based repository backend
- `managed` and `self_managed` pipelines
- Version switching (`pact use`)
- Package detection / adopt existing installs

**Phase 3 — Ecosystem**
- Central registry
- Community package contributions
- Dependency resolution
- CI/CD templates for maintainers

**Phase 4 — Advanced**
- Package signing and verification
- Rollback mechanism
- Parallel downloads and caching
- GUI client

---

## Open Questions

- **Dependencies** — how does package A declare it needs package B?
- **Conflict resolution** — what happens when two packages want the same command name?
- **Migration from Scoop/Chocolatey** — can existing manifests be ported?
- **Registry at scale** — search and index performance with thousands of packages?
- **Legal** — licensing and terms of service for a central registry?