package luautil

import (

	"github.com/yuin/gopher-lua"

)


type Capabilities struct {
    FS       FileSystem
    Registry Registry
    Env      Environment
}


type SandboxLevel int

const (
    SandboxFull    SandboxLevel = iota  // pact-runner everything 
    SandboxDry                          // pact ci no real side effects

)

func bootstrap(L *lua.LState, level SandboxLevel) *lua.LTable {
    ctx := L.NewTable()

    switch level {					// new table
    case SandboxFull:
        L.SetField(ctx, "fs",       buildFS(L))
        L.SetField(ctx, "registry", buildRegistry(L))
        L.SetField(ctx, "env",      buildEnv(L))
        L.SetField(ctx, "os",       buildOS(L))
        L.SetField(ctx, "path",     buildPath(L))

    case SandboxDry:
        L.SetField(ctx, "fs",       buildFSDry(L))
        L.SetField(ctx, "registry", buildRegistryDry(L))
        L.SetField(ctx, "env",      buildEnvDry(L))
        L.SetField(ctx, "os",       buildOS(L))       
        L.SetField(ctx, "path",     buildPath(L))     
    }

    return ctx
}





func buildFS(L *lua.LState) *lua.LTable {
    tbl := L.NewTable()
    L.SetField(tbl, "extract", L.NewFunction(fnExtract))
    L.SetField(tbl, "move",    L.NewFunction(fnMove))
    L.SetField(tbl, "remove",  L.NewFunction(fnRemove))
    L.SetField(tbl, "glob",    L.NewFunction(fnGlob))
    return tbl
}