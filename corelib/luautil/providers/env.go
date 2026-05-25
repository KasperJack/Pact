package providers


import (

	"github.com/yuin/gopher-lua"

)


type Environment interface {
    Get(L *lua.LState) int
    Set(L *lua.LState) int
    AddPath(L *lua.LState) int
}