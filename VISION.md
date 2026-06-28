# Pact Vision: Starlark-Based End-to-End Package Management

## Overview

Pact has the potential to become a fundamentally different kind of package manager for Windows by using Starlark as the universal language across the entire ecosystem — from repository automation to client-side installation.

## Current Architecture

```
Frontend CLI (runner)
    ↓
    └─→ Provides repository interface
        ↓
        Core Runtime (Starlark execution engine)
        ↓
        └─→ Calls platform abstractions
            ↓
            Platform Layer (psbridge, filesystem, registry)
            ↓
            System mutations (controlled & auditable)
```

**Key insight:** Core is functionally pure (deterministic Starlark logic) while remaining side-effect-isolated through narrow, controlled interfaces. This enables both security and auditability.

---

## Vision: Three-Layer Starlark Ecosystem

### 1. Repository Layer (Server-Side)

**Purpose:** Automated version management and package publishing.

**Concept:** Package repositories contain Starlark scripts that handle:
- Fetching latest versions from upstream (GitHub, release feeds, etc.)
- Verifying checksums and installer integrity
- Publishing to registry
- Triggering notifications

**Example implementation:**

```starlark
# In package repository (auto-runs on schedule or webhook)

def fetch_latest_version(ctx):
    """Query upstream for new releases"""
    release = ctx.http.get("https://api.github.com/repos/owner/repo/releases/latest")
    return release.tag_name.lstrip("v")

def should_update(ctx, current_version, latest_version):
    """Custom logic: skip pre-releases, respect compatibility"""
    if "rc" in latest_version or "alpha" in latest_version:
        return False
    return ctx.version_compare(latest_version, current_version) > 0

def on_new_version(ctx, version):
    """Publish new version to registry"""
    installer_url = f"https://github.com/owner/repo/releases/download/v{version}/installer.exe"
    installer = ctx.download(installer_url)
    checksum = ctx.checksum(installer, "sha256")
    ctx.publish_to_registry(version, {
        "url": installer_url,
        "checksum": checksum,
        "manifest": load_manifest(version)
    })
    ctx.notify_slack(f"Published {ctx.package_name()}@{version}")
```

**Benefits:**
- Zero-latency updates (no manual version bumps)
- Transparent, auditable update process
- Per-package custom versioning logic (semver, dates, custom schemes)
- Built-in verification before publishing

---

### 2. Package Definition Layer (Repository Metadata)

**Purpose:** Define package metadata and installation recipes.

**Concept:** Each package has a Starlark definition that:
- Describes metadata (name, dependencies, update strategy)
- Defines installation logic for any version
- Specifies version comparison rules
- Handles migrations between versions

**Example implementation:**

```starlark
# python/package.star

def metadata():
    return {
        "name": "python",
        "description": "Python interpreter",
        "maintainer": "community",
        "portable": True,
        "auto_update": True,
        "versioning": "semver",  # or custom function
    }

def install(ctx):
    """Install logic executed on client side"""
    version = ctx.version()
    installer = ctx.download(ctx.latest_url(version))
    ctx.verify_checksum(installer, ctx.checksum(version))
    ctx.extract_to(ctx.install_dir())
    ctx.add_path(ctx.install_dir() + "/bin")
    ctx.verify_install("python --version")

def migrate(ctx, from_version, to_version):
    """Handle version-specific migrations"""
    if ctx.version_parse(from_version).major < 3:
        ctx.backup_config()
        ctx.move_config("/old/path", "/new/path")
    
    if from_version < "3.11":
        ctx.run_script("upgrade_3_11.py")

def uninstall(ctx):
    """Cleanup on removal"""
    ctx.remove_path()
    ctx.cleanup_config()

def compare_versions(v1, v2):
    """Custom version comparison if needed"""
    # Default: semantic versioning
    # Can override for custom schemes (date-based, channels, etc.)
    return ctx.default_compare(v1, v2)
```

**Benefits:**
- Per-package versioning strategies (not one-size-fits-all)
- Conditional logic for migrations and upgrades
- Clear separation: metadata vs. installation logic
- Composable with other packages

---

### 3. Client Layer (Installation & Management)

**Purpose:** Execute Starlark manifests with system isolation.

**Concept:** Runner CLI provides repository interface to Core, which sandboxes execution through controlled platform abstractions.

**Current state (partial):**

```go
// runner/main.go
func installCmd(pkg, version string, arch platform.Arch) error {
    manifest := repo.Resolve(pkg, version)  // Fetch from repository
    return core.Execute(manifest, arch)      // Execute in sandbox
}
```

**What needs to be added:**

1. **Repository Interface** — Define how Core asks for packages
   ```go
   type Repository interface {
       Resolve(pkg, version string) (*Manifest, error)
       Search(query string) ([]*Package, error)
       LatestVersion(pkg string) (string, error)
   }
   ```

2. **Multiple Repository Backends:**
   - Local filesystem (development)
   - Git-based (low overhead, version control)
   - HTTP registry (central authority)
   - Hybrid (local cache + remote fallback)

3. **Expanded Starlark Standard Library:**
   - HTTP client (for fetching, API calls)
   - Cryptography (checksums, signatures)
   - Archive extraction (zip, 7z, msi, exe)
   - Process execution with timeout
   - System inspection (arch, Windows version, admin status)

---

## Implementation Roadmap

### Phase 1: Foundation (Minimal Viable Product)
- [ ] Define Repository interface in Core
- [ ] Implement local filesystem backend
- [ ] Package 3-5 critical applications (Python, Git, Node, VS Code, 7-Zip)
- [ ] Wire up basic auto-update logic
- [ ] Publish proof-of-concept packages

### Phase 2: Infrastructure
- [ ] Implement Git-based repository backend
- [ ] Add HTTP client to Starlark sandbox
- [ ] Expand standard library (crypto, archives)
- [ ] Create central registry (optional HTTP endpoint)
- [ ] Package migration tools from Scoop

### Phase 3: Ecosystem
- [ ] Community package contributions
- [ ] Dependency resolution between packages
- [ ] Lockfile format for reproducible installs
- [ ] CI/CD templates for package maintainers
- [ ] Performance optimization (parallel downloads, caching)

### Phase 4: Advanced Features
- [ ] Custom versioning schemes (date-based, channels)
- [ ] Rollback mechanism (undo installs)
- [ ] Package signing & verification
- [ ] Differential updates (delta compression)
- [ ] GUI client

---

## Architectural Advantages

### 1. Universal Language
- **Single DSL** across repository automation, package definitions, and installation
- **Lower contribution barrier** — Package maintainers learn Starlark once
- **Consistency** — Same logic runs everywhere; no translation errors

### 2. Security Through Isolation
- **Starlark sandbox** — Scripts cannot break out
- **Controlled platform layer** — All OS mutations go through auditable interfaces
- **Explicit capabilities** — Package author declares what they need (registry write, filesystem access, etc.)

### 3. Flexibility
- **Per-package customization** — Versioning, migrations, dependencies
- **Pluggable repositories** — Mix and match backends
- **Extensible by design** — New Starlark functions can be added without core changes

### 4. Reproducibility
- **Deterministic execution** — Same script + inputs = same outcome
- **Auditability** — Every operation is traceable
- **Testability** — Mock platform layer; no real side effects during testing

### 5. Comparison to Existing Managers

| Feature | Scoop | Chocolatey | Winget | **Pact** |
|---------|-------|-----------|--------|---------|
| Language | PowerShell/JSON | PowerShell/.NET | YAML (rigid) | Starlark (flexible) |
| Custom versioning | ❌ | ❌ | ❌ | ✅ |
| Conditional logic | ❌ (JSON limitation) | ✅ | ❌ | ✅ |
| Auto-update scripts | ❌ | ❌ | ❌ | ✅ |
| Sandbox | ❌ | ❌ | Limited | ✅ |
| Performance | Slow (26+ sec for search) | Slow (.NET startup) | Medium | Fast (Go + caching) |
| Cross-platform | ❌ | ❌ | ✅ | ❌ (Windows-focused) |

---

## Technical Debt & Open Questions

1. **Repository Distribution**
   - How do users discover and subscribe to repositories?
   - Decentralized (like npm) or centralized?
   - Package signing and trust model?

2. **Dependency Resolution**
   - How to handle package A depends on package B?
   - Version constraint syntax?
   - Conflict resolution?

3. **Licensing & Legal**
   - Copyright of packaged applications?
   - Terms of service for central registry?

4. **Performance at Scale**
   - How many packages can be efficiently indexed?
   - Search performance with thousands of packages?
   - Update frequency & latency?

5. **Backward Compatibility**
   - Migration path from Scoop/Chocolatey?
   - Can existing package manifests be ported?

---

## Why This Matters

Windows lacks a modern, developer-friendly package manager. Scoop is the closest thing but is limited by its JSON-based design and lacks automation. Pact has the opportunity to leapfrog by using Starlark as a universal language, enabling:

- **For users:** Fast, predictable, safe package management
- **For maintainers:** Simple, expressive package definitions with zero boilerplate
- **For developers:** A platform for building package management tools and automation

The architecture is already sound. Filling in the repository layer and packaging initial applications is the critical next step.

---

## Next Steps

1. **Design the Repository Interface** — Formalize Core's expectations
2. **Implement Local Backend** — Start with file-based packages
3. **Port 3-5 Key Packages** — Proof of concept
4. **Gather Feedback** — Early adopters on the vision
5. **Iterate on Design** — Adjust based on real-world usage
