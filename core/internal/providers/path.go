package providers

import (
	"fmt"
	"path/filepath"

	"go.starlark.net/starlark"
    stdos "os"
)

type Path interface {
	Join(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error)
    MoveAll(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error)
}

type path struct{}

func NewTestPath() Path {

	return &path{}
}


func (p *path) Join(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
    if len(kwargs) > 0 {
        return nil, fmt.Errorf("%s: unexpected keyword arguments", fn.Name())
    }

    parts := make([]string, len(args))
    for i, arg := range args {
        s, ok := starlark.AsString(arg)
        if !ok {
            return nil, fmt.Errorf("%s: for parameter %d: got %s, want string", fn.Name(), i+1, arg.Type())
        }
        parts[i] = s
    }

    result := filepath.Join(parts...)
    return starlark.String(result), nil
}







func (f *path) MoveAll(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
    var src, dst starlark.String

    if err := starlark.UnpackArgs(fn.Name(), args, kwargs, "src", &src, "dst", &dst); err != nil {
        return nil, err
    }

    
    // check  src exist and is a dir 
    srcBase := filepath.Dir(string(src))

    srcInfo, err := stdos.Stat(srcBase)
    if err != nil {
        if stdos.IsNotExist(err) {
            return nil, fmt.Errorf("%s: source %q does not exist", fn.Name(), srcBase)
        }
        return nil, fmt.Errorf("%s: %w", fn.Name(), err)
    }
    if !srcInfo.IsDir() {
        return nil, fmt.Errorf("%s: source %q is not a directory", fn.Name(), srcBase)
    }





    // check dst exists and is a dir
    dstInfo, err := stdos.Stat(string(dst))
    if err != nil {
        if stdos.IsNotExist(err) {
            return nil, fmt.Errorf("%s: destination %q does not exist", fn.Name(), string(dst))
        }
        return nil, fmt.Errorf("%s: %w", fn.Name(), err)
    }
    if !dstInfo.IsDir() {
        return nil, fmt.Errorf("%s: destination %q is not a directory", fn.Name(), string(dst))
    }



    matches, err := filepath.Glob(string(src))
    if err != nil {
        // Only returned for a malformed pattern (filepath.ErrBadPattern).
        return nil, fmt.Errorf("%s: %w", fn.Name(), err)
    }



    for _, m := range matches {
        target := filepath.Join(string(dst), filepath.Base(m))
        if err := stdos.Rename(m, target); err != nil {
            return nil, fmt.Errorf("%s: moving %s: %w", fn.Name(), m, err)
        }
    }

    return starlark.None, nil
}