# pact

Pact is a software distribution platform for Windows. Packages are defined 
in Starlark — a safe, sandboxed scripting language.

> ⚠️ Early stage

## What's here

- **Runtime** — executes Starlark package manifests
- **Platform bridge** — Windows (PowerShell)
- **CI runner** — runs manifests in a pipeline
- **Local runner** — test and debug manifests locally

## Package manifests

```python
def install():
    path.add(install_dir)
    shortcut.create("My App", install_dir + "/app.exe")
```