package providers

import (

	"github.com/yuin/gopher-lua"

)


type FileSystem interface {
    Extract(L *lua.LState) int
    Move(L *lua.LState) int
    Remove(L *lua.LState) int
    Glob(L *lua.LState) int
}