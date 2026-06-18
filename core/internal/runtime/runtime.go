package runtime

import (
    "fmt"
    //"go.starlark.net/starlark"
	"go.starlark.net/starlark"
    "go.starlark.net/syntax"
)



type installContext struct {
    thread  *starlark.Thread
    globals starlark.StringDict
}


func NewIinstallContext(script []byte) (*installContext,error) {

    thread := &starlark.Thread{Name: "main"}
    opts := syntax.FileOptions{}

    globals, err := starlark.ExecFileOptions(&opts, thread, "script.star", script, nil)  //top level eval
    if err != nil {
        return nil, err
    }

    return &installContext{thread: thread, globals: globals}, nil

}



func (ctx *installContext) Run() error {
    result, ok := ctx.globals["result"]
    if !ok {
        return fmt.Errorf("result not found in script")
    }
    // starlark.Call(ctx.thread, fn, starlark.Tuple(args), nil)
    fmt.Println(result)
    return nil
}