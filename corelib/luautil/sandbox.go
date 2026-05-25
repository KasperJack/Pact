package luautil

import (

	"github.com/yuin/gopher-lua"
	"Pact/corelib/luautil/providers"

)


type Capabilities struct {
    FS       providers.FileSystem
    Registry providers.Registry
    Env      providers.Environment
}


type SandboxLevel int

const (
    SandboxFull    SandboxLevel = iota  // pact-runner everything 
    SandboxDry                          // pact ci no real side effects

)

func bootstrap(L *lua.LState, caps Capabilities) *lua.LTable {
    ctx := L.NewTable()
    L.SetField(ctx, "fs",       buildFS(L, caps.FS))
    L.SetField(ctx, "reg", buildRegistry(L, caps.Registry))
    L.SetField(ctx, "env",      buildEnv(L, caps.Env))
    return ctx
}





func buildFS(L *lua.LState, fs providers.FileSystem) *lua.LTable {
    tbl := L.NewTable()
    L.SetField(tbl, "extract", L.NewFunction(fs.Extract))
    L.SetField(tbl, "move",    L.NewFunction(fs.Move))
    L.SetField(tbl, "remove",  L.NewFunction(fs.Remove))
    L.SetField(tbl, "glob",    L.NewFunction(fs.Glob))
    return tbl
}



func buildRegistry(L *lua.LState, reg providers.Registry) *lua.LTable {
    tbl := L.NewTable()
    L.SetField(tbl, "read", L.NewFunction(reg.Read))
    L.SetField(tbl, "write",    L.NewFunction(reg.Write))
    L.SetField(tbl, "delete",  L.NewFunction(reg.Delete))

    return tbl
}



func buildEnv(L *lua.LState, env providers.Environment) *lua.LTable {
    tbl := L.NewTable()
    L.SetField(tbl, "read", L.NewFunction(env.AddPath))
    L.SetField(tbl, "write",    L.NewFunction(env.Get))
    L.SetField(tbl, "delete",  L.NewFunction(env.Set))

    return tbl
}