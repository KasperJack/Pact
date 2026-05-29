package luautil

import (
	"Pact/corelib/model"
	"fmt"
    //"strings"
    "log"
	"github.com/yuin/gopher-lua"
    //"github.com/yuin/gopher-lua/parse"
    //"github.com/yuin/gopher-lua/ast"
)

type PackageEvalContext struct {
    l   *lua.LState
    pkg *model.Package
}

func NewPackageEvalContext(pkg *model.Package) *PackageEvalContext {
    ctx := &PackageEvalContext{
        l:   lua.NewState(lua.Options{SkipOpenLibs: true}), // no open libs
        pkg: pkg,
    }
    ctx.l.SetGlobal("package", ctx.l.NewFunction(ctx.fnPackage))
    ctx.l.SetGlobal("printFromGo", ctx.l.NewFunction(ctx.fnPrintFromGo))
    return ctx
}

func (ctx *PackageEvalContext) fnPackage(L *lua.LState) int {
    tbl := L.CheckTable(1)

	checkNoExtraKeys(L,tbl,allowedKeysFromStruct(&model.Package{})) // rasies a lua error 

    ctx.pkg.PackageIdentifier = requiredString(L, tbl, "package_identifier")
    ctx.pkg.Name              = requiredString(L, tbl, "name")
    ctx.pkg.Versioning        = requiredString(L, tbl, "versioning",)
    ctx.pkg.Description       = optionalString(L, tbl,"description","deflt")
    ctx.pkg.Homepage          = optionalString(L, tbl,"homepage","deflt")
    ctx.pkg.License           = optionalString(L, tbl,"license","deflt")


    installFn, ok := L.GetField(tbl, "install").(*lua.LFunction)
    if ok {
        ctx.pkg.InstallFn = installFn

    }




    return 0
}

func (ctx *PackageEvalContext) Eval(luaData []byte) error {
    return ctx.l.DoString(string(luaData))
    //return parceLua(luaData)
}     

func (ctx *PackageEvalContext) Close() {
    ctx.l.Close()
}




func (ctx *PackageEvalContext) fnPrintFromGo(l *lua.LState) int {
    arg := l.CheckString(1) // get the first argument from Lua
    fmt.Println(arg)
    return 0 // number of return values pushed onto the stack
}


func (ctx *PackageEvalContext) RunInstall() error {
    if ctx.pkg.InstallFn == nil {
        return fmt.Errorf("no install function defined")
    }
    return ctx.l.CallByParam(lua.P{
        Fn:      ctx.pkg.InstallFn,
        NRet:    0,
        Protect: true,
    },bootstrap(ctx.l))
}

















type TestEvalContext struct {
    l   *lua.LState

}


func NewTestEvalContext() *TestEvalContext {
    ctx := &TestEvalContext{
        l:   lua.NewState(lua.Options{SkipOpenLibs: true}), // no open libs
    }
    //ctx.l.SetGlobal("package", ctx.l.NewFunction(ctx.fnPackage))
    //ctx.l.SetGlobal("printFromGo", ctx.l.NewFunction(ctx.fnPrintFromGo))
    return ctx
}




func (ctx *TestEvalContext) Eval(luaData []byte) error {
    return ctx.l.DoString(string(luaData))
    //return parceLua(luaData)
}     

func (ctx *TestEvalContext) Close() {
    ctx.l.Close()
}


func (ctx *TestEvalContext) RunInstall() error {


    onInstallFunc := ctx.l.GetGlobal("install")

	if onInstallFunc.Type() != lua.LTFunction {
		log.Fatal("install function not found")
	}
    ctx.l.SetFEnv(onInstallFunc, bootstrap(ctx.l))


    err := ctx.l.CallByParam(lua.P{
		Fn:      onInstallFunc,
		NRet:    0,
		Protect: true,
	})
	if err != nil {
		return err
	}

    return nil

}