package providers

import (
	//"fmt"
	"runtime"
	"go.starlark.net/starlark"
)

type Os interface {


	GArch() starlark.String



	IsX64(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error)
	IsX86(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error)
	IsArm64(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error)
}





type os struct{}

func NewOs() Os {
	return &os{}
}



func (o *os) IsX64(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := starlark.UnpackArgs(fn.Name(), args, kwargs); err != nil {
		return nil, err
	}
	return starlark.Bool(runtime.GOARCH == "amd64"), nil
}

func (o *os) IsX86(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := starlark.UnpackArgs(fn.Name(), args, kwargs); err != nil {
		return nil, err
	}
	return starlark.Bool(runtime.GOARCH == "386"), nil
}

func (o *os) IsArm64(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := starlark.UnpackArgs(fn.Name(), args, kwargs); err != nil {
		return nil, err
	}
	return starlark.Bool(runtime.GOARCH == "arm64"), nil
}

func (o *os)GArch() starlark.String {

	return starlark.String("x64")
}