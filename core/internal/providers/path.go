package providers

import (
	//"fmt"
	"path/filepath"

	"go.starlark.net/starlark"
)

type Path interface {
	Join(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error)
}

type path struct{}

func NewTestPath() Path {

	return &path{}
}


func (p *path) Join(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
    var a, b starlark.String

    if err := starlark.UnpackArgs(fn.Name(), args, kwargs, "a", &a, "b", &b); err != nil {
        return nil, err
    }

    result := filepath.Join(string(a), string(b))
    return starlark.String(result), nil 
}



