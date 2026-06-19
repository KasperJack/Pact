package runtime
import (
    //"fmt"
    //"go.starlark.net/starlark"
	"go.starlark.net/starlark"
    //"go.starlark.net/syntax"
	"github.com/kasperjack/pact/core/internal/providers"
	"go.starlark.net/starlarkstruct"
)

type Capabilities struct {
    Path       providers.Path
    //Registry providers.Registry
    //Env      providers.Environment
}

func buildPath(Path providers.Path) starlark.StringDict {
    return starlark.StringDict{
        "join": starlark.NewBuiltin("join", Path.Join),

    }
}


func newPreDeclared(caps Capabilities) starlark.StringDict {


	predeclared := starlark.StringDict{


        "path":  starlarkstruct.FromStringDict(starlarkstruct.Default, buildPath(caps.Path)),
        //"reg": starlarkstruct.FromStringDict(starlarkstruct.Default, buildRegistry(caps.Registry)),
        //"env": starlarkstruct.FromStringDict(starlarkstruct.Default, buildEnv(caps.Env)),
    }





	return predeclared
}