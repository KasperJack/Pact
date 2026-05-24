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
    return ctx
}

func (ctx *PackageEvalContext) fnPackage(L *lua.LState) int {
    tbl := L.CheckTable(1)

	checkNoExtraKeys(L,tbl,allowedKeysFromStruct(&model.Package{})) // rasies an lua error 

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
    //return parceLua(luaData)
}     

func (ctx *PackageEvalContext) Close() {
    ctx.l.Close()
}


func parceLua (LuaData []byte) error {
    chunk, err := parse.Parse(strings.NewReader(string(LuaData)), "<string>")
    if err != nil {
        return  err
    }

    findCalls(chunk)
    return nil
}




















func findCalls(stmts []ast.Stmt) {
    for _, stmt := range stmts {
        walkStmt(stmt)
    }
}

func walkStmt(stmt ast.Stmt) {
    switch s := stmt.(type) {
    case *ast.FuncCallStmt:
        walkExpr(s.Expr)
    case *ast.FuncDefStmt:
        findCalls(s.Func.Stmts)
    case *ast.IfStmt:
        findCalls(s.Then)
        findCalls(s.Else)
    case *ast.ReturnStmt:
        for _, e := range s.Exprs {
            walkExpr(e)
        }
    // add more cases as needed...
    }
}

func walkExpr(expr ast.Expr) {
    switch e := expr.(type) {
    case *ast.FuncCallExpr:
        fmt.Printf("called: %v, args: %d\n", e.Func, len(e.Args))
        // recurse into args too
        for _, arg := range e.Args {
            walkExpr(arg)
        }
    }
}