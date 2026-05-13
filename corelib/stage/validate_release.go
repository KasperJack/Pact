package stage

import (
	"fmt"
    "Pact/corelib/client"

)

type Release struct {
    PackageIdentifier string `hcl:"package_identifier"`
	Url    string `hcl:"url"`
	Version string `hcl:"version"`
    Hash string `hcl:"hash"`
}

func (r *Release) Validate (mode string) error {

    var s client.RepositorySource

    switch mode {

    case "l":
        s = client.NewLFilesystemSource("/somewhere")
    case "r":
        s,_ = client.NewGithubSource("kasperjack/pact","main")
    default:
        panic("ass")

    }

    rc := client.NewRepoClient(s)
    ok := rc.PackageExists(r.PackageIdentifier)

    fmt.Println("validating release")
    return nil
}  
