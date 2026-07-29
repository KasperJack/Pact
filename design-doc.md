# pact — Design Doc 29/7/26

Status: **living document**. Treat this as the source of truth for architecture
decisions. When something in code disagrees with this doc, either the code is
wrong or the doc is stale — fix whichever one is actually wrong, don't let
them silently diverge.

---

## 1. Goals

- Install, list, and manage **portable** Windows apps and CLI tools (no
  installers-with-side-effects; no registry pollution as a hard requirement).
- Support installing **any historical version** of a package, not just latest.
- Separate "what upstream released" from "how we package it" so packaging
  bugs can be fixed without needing a new upstream release.
- Support **multiple architectures** (x64, arm64, x86, and arch-independent)
  per package, where install steps can genuinely differ per arch.
- Stay simple enough to hand-write manifests for the current phase. Tooling
  gets built only once real pain is felt, not speculatively.

## 2. Non-goals (for now)

- No dependency resolution between packages (each package installs standalone).
- No epoch concept (see §4.4) unless a real upstream forces it.
- No auto-generated manifests / scrapers yet. Everything is hand-authored
  until the manual process actually hurts.
- No package removal/rollback transaction log yet — out of scope until basic
  install/list/upgrade is solid.

---

## 3. Core Versioning Model

Borrowed from Debian's `epoch:upstream_version-debian_revision` scheme, with
epoch dropped for now.

### 3.1 Two distinct identities

| Term | Example | Meaning |
|---|---|---|
| `UpstreamVersion` | `2.7.0` | Exactly what upstream calls their release. Copied verbatim when legal; otherwise mangled once, by hand, at manifest-authoring time (strip `_`, leading `v`, etc.) |
| `FullVersion` | `2.7.0-2` | `UpstreamVersion` + `-Revision`. This is **our** packaging identity — bumped whenever the manifest changes but the upstream binary doesn't. |
| `Revision` | `2` (int) | How many times we've re-packaged this exact upstream version. Starts at `1`, never `0`. |

Rule: **the full manifest, per revision, is a standalone, complete file — never a diff/patch against the previous revision.** A revision bump means "copy the previous manifest, fix what was wrong, write it as a new file." The old revision's manifest is never mutated in place.

### 3.2 Revision counters are per-architecture, independent

Different architectures do **not** share a revision counter. If x64's
manifest has a bug and arm64's doesn't, only x64 bumps:

```
releases/x64/windirstat/2.7.0-2/    ← real fix
releases/arm64/windirstat/2.7.0-1/  ← untouched, no reason to bump
```

Rationale: avoids meaningless copy-forward manifests on unaffected arches.
The cost (per-arch revision numbers can look "out of sync") is acceptable
because most users only ever see the collapsed `UpstreamVersion` view (§6.1)
— the raw revision mismatch is rarely user-visible.

### 3.3 What a revision bump can and cannot change

Allowed without a revision bump: nothing. Any manifest edit = new revision.
(Simplicity rule: revision count should always equal actual manifest edit
count for that arch.)

Typical reasons to bump revision:
- Wrong install flags / silent-install args
- Wrong extraction path
- Wrong or stale checksum for the *same* upstream artifact
- Any other packaging-only fix

Reasons to bump `UpstreamVersion` instead (new folder, revision resets to 1):
- Upstream actually shipped new code/binary

### 3.4 Epoch — deferred, not deleted

Epoch exists in Debian to handle upstream changing its versioning scheme in
a way that breaks natural sort order (e.g. date-based `20230401` → semver
`1.0.0`). We are **not** implementing epoch parsing yet. If this scenario
is hit for a real package, the interim fix is a manifest-level
`sort_override` escape hatch rather than building full epoch support
pre-emptively. Revisit if this happens more than once.

---

## 4. Repository Layout (on disk)

```
packages/
  <pkg>/
    package.hcl              ← identity metadata, version-independent

releases/
  <arch>/                    ← x64 | arm64 | x86 | noarch
    <pkg>/
      index.hcl              ← per-arch version index (hand-written for now)
      <full_version>/
        package.hcl           ← standalone install recipe for this exact revision
```

Notes:
- `<arch>` is a **directory boundary**, not just a manifest field — because
  install steps (not just the download URL) can differ per architecture.
- `noarch` is reserved for packages with no compiled binary (scripts, jars,
  anything identical across all arches) — avoids triplicating identical
  manifests across x64/arm64/x86.
- `index.hcl` is per-arch. There is no single global "latest version" for a
  package across all arches simultaneously — only per-arch latest.

---

## 5. Manifest Formats

### 5.1 `packages/<pkg>/package.hcl` — identity, version-independent

```hcl
package     = "windirstat"
description = "Disk usage statistics viewer"
homepage    = "https://windirstat.net"
license     = "GPL-2.0"
```

Anything true regardless of version or arch lives here. Never duplicated
per-release.

### 5.2 `releases/<arch>/<pkg>/<full_version>/package.hcl` — install recipe

```hcl
package          = "windirstat"
version          = "2.7.0-1"        # FullVersion
upstream_version = "2.7.0"
revision         = 1
url              = "https://github.com/windirstat/windirstat/releases/download/release%2Fv2.7.0/WinDirStat.zip"
sha256           = "3aad34ab829ddbfc7f859c22554981b19b1b04dd4cb87b643901e4eb9e2d85c4"
size_mb          = 8
architecture     = "amd64"

# install behavior — fields TBD as install logic gets built, but must live
# here (self-contained recipe principle). placeholder shape:
# install_type = "zip_extract" | "installer_exe" | "msi"
# extract_to   = "%LOCALAPPDATA%\\pkgmgr\\windirstat"
# post_install = [...]
```

Rule: this file must contain everything needed to install this exact
revision with zero reference to any other file (other than the parent
`package.hcl` for display metadata like description/homepage).

### 5.3 `releases/<arch>/<pkg>/index.hcl` — per-arch version index

```hcl
latest_version = "2.7.0-2"

version_mappings = {
  "2.6.2" = "2.6.2-1"
  "2.7.0" = "2.7.0-2"
  "2.6.0" = "2.6.0-2"
  "2.5.0" = "2.5.0-1"
}

# optional, populated only when a revision is known-broken:
# yanked = {
#   "2.7.0-1" = "silent-install flag caused leftover UI popup"
# }
```

Hand-written for now. `version_mappings` maps `UpstreamVersion → latest good
FullVersion for that upstream version`. Consumers should treat this file,
not directory scanning, as the source of truth while it's hand-maintained —
revisit once/if a generator replaces hand-editing.

---

## 6. Resolution Rules

### 6.1 Default `list` — collapsed to upstream versions

```
> pkgmgr list windirstat
2.6.0
2.6.2
2.7.0   (latest)
```

### 6.2 Expanded `list --all` — full per-arch revision history

```
> pkgmgr list windirstat --all --arch x64
2.6.0-1
2.6.0-2
2.6.2-1
2.7.0-1
2.7.0-2   (latest)
```

### 6.3 Install resolution

| Command | Resolves to |
|---|---|
| `pkgmgr install windirstat` | latest `UpstreamVersion`, latest `Revision` for detected arch |
| `pkgmgr install windirstat@2.7.0` | latest `Revision` of that `UpstreamVersion` (via `version_mappings`) |
| `pkgmgr install windirstat@2.7.0-1` | exact pin — installs precisely that manifest, bugs included |

Pinning to a yanked revision: **warn, require confirmation, still allow it.**
Never silently substitute the fixed revision — that breaks the reproducibility
guarantee that pinning exists for.

### 6.4 Upgrade

`pkgmgr upgrade <pkg>` compares the installed `FullVersion` against the
current arch's `index.hcl` `latest_version` using the version-compare
algorithm (§7). If newer, re-runs install steps per the new manifest.
Whether this is full uninstall/reinstall or a lighter in-place step is an
implementation detail, not a versioning-scheme concern.

---

## 7. Version Comparison Algorithm

Use dpkg's `verrevcmp` algorithm (no epoch stage, since epoch is deferred
per §3.4):

1. Compare `UpstreamVersion` strings using alternating digit/non-digit run
   comparison, where digit runs compare numerically and `~` sorts before
   everything (including empty string) if we ever need pre-release tags.
2. If equal, compare `Revision` as plain integers.

Decision: **use an existing Go library** rather than hand-rolling this —
`github.com/knqyf263/go-deb-version` or `pault.ag/go/debian/version`. Both
implement the real dpkg algorithm faithfully; no reason to maintain a
parallel implementation for something this well-specified. If epoch is
never needed, the epoch component can just always be `0` when using either
library.

---

## 8. Local State (installed packages)

Not yet fully designed — placeholder shape based on what's been discussed:

```json
{
  "name": "windirstat",
  "architecture": "x64",
  "installed_version": "2.7.0-2",
  "upstream_version": "2.7.0",
  "revision": 2,
  "install_path": "C:\\Users\\...\\pkgmgr\\windirstat"
}
```

Open question, not yet resolved: where does this live (single JSON/SQLite
db vs one file per installed package)? Revisit when install logic is
actually being built.

---

## 9. Explicit Open Questions

Things intentionally left unresolved — do not silently assume an answer
elsewhere in code without updating this section:

- [ ] Exact shape of install-behavior fields in the release manifest
      (`install_type`, `post_install`, shortcut creation, etc.)
- [ ] Local install-state storage format
- [ ] Whether `sort_override` (epoch escape hatch) is ever actually needed
- [ ] Uninstall / rollback flow
- [ ] Whether `index.hcl` should eventually be generated from directory
      scanning once enough packages exist to make hand-editing painful

---

## 10. Explicitly Deferred (do not build yet)

- Manifest generator / scraper tooling
- Epoch support
- Dependency resolution between packages
- Any GUI — CLI only for now