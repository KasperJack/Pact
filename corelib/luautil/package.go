package luautil

import (
	"Pact/corelib/model"
	"fmt"
    "strings"
	"github.com/yuin/gopher-lua"
    "github.com/yuin/gopher-lua/parse"
    "github.com/yuin/gopher-lua/ast"
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
	//fmt.Println(allowedKeysFromStruct(&model.Package{}))
    return ctx
}

func (ctx *PackageEvalContext) fnPackage(L *lua.LState) int {
    tbl := L.CheckTable(1)

	checkNoExtraKeys(L,tbl,allowedKeysFromStruct(&model.Package{}))

    ctx.pkg.PackageIdentifier = requiredString(L, tbl, "package_identifier")
    ctx.pkg.Name              = requiredString(L, tbl, "name")
    ctx.pkg.Versioning        = requiredString(L, tbl, "versioning",)
    ctx.pkg.Description       = optionalString(L, tbl,"description","deflt")
    ctx.pkg.Homepage          = optionalString(L, tbl,"homepage","deflt")
    ctx.pkg.License           = optionalString(L, tbl,"license","deflt")

    return 0
}

func (ctx *PackageEvalContext) Eval(luaData []byte) error {
    //return ctx.l.DoString(string(luaData))
    return parceLua(luaData)
}     

func (ctx *PackageEvalContext) Close() {
    ctx.l.Close()
}


func parceLua (LuaData []byte) error {
    chunk, err := parse.Parse(strings.NewReader(string(LuaData)), "<string>")
    if err != nil {
        return  err
    }

    for _, stmt := range chunk {
    if fs, ok := stmt.(*ast.FuncDefStmt); ok {
        if ident, ok := fs.Name.Func.(*ast.IdentExpr); ok {
            fmt.Println(ident.Value)
            }
        }
    }

    return nil
}