# pact-runner

a local package script runner for testing and verifying pact package definitions before submitting them to the registry.

## what it does

loads a `.lua` package definition, spins up a sandboxed Lua runtime, and executes the install lifecycle against your local machine. lets you write and test a package script without needing the full pact package manager.

## what it is not

- not a package manager, it does not track installs, manage versions, or maintain a database
- not a sandbox validator, it runs real installs on your real machine
- not for end users, its a developer tool 

## workflow

```
write package.lua  →  pact-runner package.lua  →  see if it works  →  fix  →  repeat
```

## lifecycle

the runner executes hooks in this order:

```
pre_install   →  setup, registry patching, arch detection
install       →  extract, move, msi expansion
post_install  →  shortcuts, env vars, path registration
```

if any step fails the runner stops and reports exactly which hook failed and why.

## ctx

every lifecycle function receives a `ctx` object that is the entire API surface the package author has access to. nothing outside of `ctx` can interact with the system.

```lua
install = function(ctx)
    ctx.extract(ctx.dist(), ctx.staging())
    ctx.move(ctx.path.join(ctx.staging(), "x64/*"), ctx.dir())
end
```

## sandboxing

the Lua runtime boots with no standard library. the package author can only call functions explicitly registered by pact-runner. there is no `io`, no `os`, no `require`. only `ctx`.

## usage

```bash
pact-runner <package.lua>

# examples
pact-runner windirstat.lua
pact-runner python.lua
```

## output

```
installing windirstat 2.6.1
  pre_install
    ✓ detect architecture
    ✓ patch registry files
  install
    ✓ extract archive
    ✓ move x64 files
  post_install
    ✓ create shortcut
done in 1.2s
```