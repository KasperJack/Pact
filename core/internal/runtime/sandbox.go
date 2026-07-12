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
    Os providers.Os
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


func buildOs(o providers.Os) starlark.StringDict {

	os := starlark.StringDict{}


	if o == nil {

		os["is_x64"] = blocked("is_x64")
        os["is_x86"] = blocked("is_x86")
        os["is_arm64"] = blocked("is_arm64")



	}else{
        os["arch"] = o.GArch()
		os["is_x64"] = starlark.NewBuiltin("is_x64", o.IsX64)
        os["is_x86"] = starlark.NewBuiltin("is_x86", o.IsX86)
        os["is_arm64"] = starlark.NewBuiltin("is_arm64", o.IsArm64)
	}


	return os
}