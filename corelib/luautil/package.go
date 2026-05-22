package luautil

import (
	"github.com/yuin/gopher-lua"
	"Pact/corelib/model"
)

type PackageEvalContext struct {
    l   *lua.LState
    pkg *model.Package
}

func NewPackageEvalContext(pkg *model.Package) *PackageEvalContext {
    ctx := &PackageEvalContext{
        l:   lua.NewState(),
        pkg: pkg,
    }
    ctx.l.SetGlobal("package", ctx.l.NewFunction(ctx.fnPackage))
    return ctx
}

func (ctx *PackageEvalContext) fnPackage(L *lua.LState) int {
    tbl := L.CheckTable(1)

    ctx.pkg.PackageIdentifier = requiredString(L, tbl, "package_identifier")
    ctx.pkg.Name              = requiredString(L, tbl, "name")
    ctx.pkg.Versioning        = requiredString(L, tbl, "versioning",)
    ctx.pkg.Description       = optionalString(L, tbl,"description","deflt")
    ctx.pkg.Homepage          = optionalString(L, tbl,"homepage","deflt")
    ctx.pkg.License           = optionalString(L, tbl,"license","deflt")

    return 0
}

func (ctx *PackageEvalContext) Eval(luaData []byte) error {
    return ctx.l.DoString(string(luaData))
}

func (ctx *PackageEvalContext) Close() {
    ctx.l.Close()
}