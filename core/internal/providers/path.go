package providers

import (
	//"fmt"
	"fmt"

	"go.starlark.net/starlark"
)

type Path interface {
	Join(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error)
}

type path struct{}

func NewTestPath() Path {

	return &path{}
}


func (*path) Join(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error){

	fmt.Println("patttttt")

	return starlark.None,nil
}