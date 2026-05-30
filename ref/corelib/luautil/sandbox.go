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
// caps Capabilities
func bootstrap(L *lua.LState) *lua.LTable {

    caps := Capabilities{
        FS: new(providers.TestFs),
        // the rest
    }
    ctx := L.NewTable()
    L.SetField(ctx, "fs",       buildFS(L, caps.FS))
    //L.SetField(ctx, "reg", buildRegistry(L, caps.Registry))
    //L.SetField(ctx, "env",      buildEnv(L, caps.Env))
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


func buildOs(L *lua.LState) *lua.LTable {
    tbl := L.NewTable()
    L.SetField(tbl, "x64",   lua.LBool(runtime.GOARCH == "amd64"))
    L.SetField(tbl, "x86",   lua.LBool(runtime.GOARCH == "386"))
    L.SetField(tbl, "arm64", lua.LBool(runtime.GOARCH == "arm64"))
    return tbl
}