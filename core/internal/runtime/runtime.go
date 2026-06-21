package runtime

import (
	"fmt"
	//"go.starlark.net/starlark"
	"github.com/kasperjack/pact/core/internal/providers"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"go.starlark.net/syntax"
)



type installContext struct {
    thread  *starlark.Thread
    globals starlark.StringDict
    predeclared starlark.StringDict
}





func NewIinstallContext(script []byte) (*installContext,error) {

    
    thread := &starlark.Thread{Name: "main"}

    opts := syntax.FileOptions{
    TopLevelControl: true,
    GlobalReassign:  false, // default, but be explicit
    }



    predeclared := starlark.StringDict{
        "path":  starlarkstruct.FromStringDict(starlarkstruct.Default, buildPath(nil)),
        //"reg": starlarkstruct.FromStringDict(starlarkstruct.Default, buildRegistry(caps.Registry)),
        //"env": starlarkstruct.FromStringDict(starlarkstruct.Default, buildEnv(caps.Env)),
        "install_dir": starlark.String("/usr/local"),
        "version":     starlark.String("1.0.0"),
    }

    

    // goroutine needed here with a timeout 
    // inject globals and functions                                               //predeclared                         
    globals, err := starlark.ExecFileOptions(&opts, thread, "script.star", script, predeclared)  //top level eval
    if err != nil {
        return nil, err
    }


    for name := range predeclared {
        if _, ok := globals[name]; ok {
        return nil, fmt.Errorf("%s is a reserved name and cannot be reassigned", name)
        }
    }



    return &installContext{thread: thread, globals: globals,predeclared: predeclared}, nil

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

    p:= providers.NewTestPath()

    ctx.predeclared["path"] = starlarkstruct.FromStringDict(starlarkstruct.Default, buildPath(p))

    _, err := starlark.Call(ctx.thread, callable, nil, nil)
    if err != nil {
        return err
    }

    return nil
}



















func logBuiltin (thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {

        var action starlark.String

        if err := starlark.UnpackArgs(fn.Name(), args, kwargs, "action", &action); err != nil {
        return nil, err
    }

        fmt.Printf("log:%s \n",action.GoString())
        return starlark.None, nil
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

