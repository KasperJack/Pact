package providers

import (

	"github.com/yuin/gopher-lua"
    "fmt"

)


type FileSystem interface {
    Extract(L *lua.LState) int
    Move(L *lua.LState) int
    Remove(L *lua.LState) int
    Glob(L *lua.LState) int
}


type TestFs struct {}



func (* TestFs) Extract(L *lua.LState)int {
    src := L.CheckString(1)
	dst := L.CheckString(2)


	fmt.Printf("extracting %s to %s \n",src,dst)

	return 0
}

func (* TestFs) Move(L *lua.LState)int {
    src := L.CheckString(1)
	dst := L.CheckString(2)


	fmt.Printf("moving %s to %s \n",src,dst)

	return 0
}


func (* TestFs) Remove(L *lua.LState)int {
    target := L.CheckString(1)



	fmt.Printf("removeing %s \n",target)

	return 0
}


func (* TestFs) Glob(L *lua.LState)int {
  

	fmt.Println("runnning Glob")

	return 0
}