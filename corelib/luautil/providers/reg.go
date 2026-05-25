package providers

import (

	"github.com/yuin/gopher-lua"

)


type Registry interface {
    Read(L *lua.LState) int
    Write(L *lua.LState) int
    Delete(L *lua.LState) int
}