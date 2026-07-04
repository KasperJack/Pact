package runtime
import (
    "fmt"
    //"go.starlark.net/starlark"
	"go.starlark.net/starlark"
    //"go.starlark.net/syntax"
	"github.com/kasperjack/pact/core/internal/providers"
	//"go.starlark.net/starlarkstruct"
)

type Capabilities struct {
    Path       providers.Path
    //Registry providers.Registry
    //Env      providers.Environment
}

func builtin(name string, fn func(args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error)) *starlark.Builtin {
    return starlark.NewBuiltin(name, func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
        return fn(args, kwargs)
    })
}

func blocked(name string) *starlark.Builtin {
    return builtin(name, func(args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
        return nil, fmt.Errorf("%s can only be called inside install", name) // add line number 
    })
}


func buildPath(p providers.Path) starlark.StringDict {

	path := starlark.StringDict{}


	if p == nil {

		path["join"] = blocked("join")



	}else{

		path["join"] = starlark.NewBuiltin("join", p.Join)

	}


	return path
}


