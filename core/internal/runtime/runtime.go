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
    opts := syntax.FileOptions{TopLevelControl: true}



    predeclared := starlark.StringDict{
    "print_from_go": starlark.NewBuiltin("print_from_go", PrintFromGo ),
    }

    

    // goroutine needed here with a timeout 
    // inject globals and functions                                               //predeclared                         
    globals, err := starlark.ExecFileOptions(&opts, thread, "script.star", script, predeclared)  //top level eval
    if err != nil {
        return nil, err
    }

    return &installContext{thread: thread, globals: globals}, nil

}



func (ctx *installContext) Run() error {
    fn, ok := ctx.globals["install"]
    if !ok {
        return fmt.Errorf("no install func was defined")
    }

    callable, ok := fn.(starlark.Callable)
    if !ok {
        return fmt.Errorf("install is not callable")
    }

    if callable.(*starlark.Function).NumParams() != 0 {
        return fmt.Errorf("install must take no arguments, but got %d", callable.(*starlark.Function).NumParams())
    }



    _, err := starlark.Call(ctx.thread, callable, nil, nil)
    if err != nil {
        return err
    }

    return nil
}




func PrintFromGo (thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
    /*
    var name starlark.String
    var age starlark.String
    
    if err := starlark.UnpackArgs(fn.Name(), args, kwargs,"name", &name,"age",&age,); err != nil {
        return nil, err
    }
    
    fmt.Println(name.GoString())
    fmt.Println(age.GoString())
*/
    fmt.Println("age.GoString()")
    return starlark.None, nil
    }