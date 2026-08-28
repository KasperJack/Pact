# Pact Script Reference

## Automatic Behaviors
Pact handles the following before `install()` is called:
- Downloads the archive matching the target arch and extracts it into the staging dir
- Falls back to x86 if no x64 binary is available on an x64 machine
- Respects arch overrides passed via the CLI (e.g. --arch x86 on an x64 machine)

## pkg — Package Info

| Variable        | Type | Example                              | Description                          |
|-----------------|------|--------------------------------------|--------------------------------------|
| pkg.name        | str  | windirstat                           | Package name                         |
| pkg.version     | str  | 2.6.4                                | Installed version                    |
| pkg.arch        | str  | x64                                  | Target arch pact resolved (x86, x64, arm64) |
| pkg.staging     | str  | C:/pact/staging/windirstat/2.6.4/    | Extracted archive contents land here |
| pkg.install_dir | str  | C:/pact/packages/windirstat/2.6.4/   | Final install destination            |

## os — Host System

| Function      | Returns | Description                        |
|---------------|---------|------------------------------------|
| os.arch()     | str     | Actual OS arch (x86, x64, arm64)   |
| os.is_x64()   | bool    | True if OS arch is x64             |
| os.is_x86()   | bool    | True if OS arch is x86             |
| os.is_arm64() | bool    | True if OS arch is arm64           |

> `pkg.arch` is what pact chose to install (may differ from `os.arch()` due to availability or CLI override)