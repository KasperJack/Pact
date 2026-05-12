package stage

import (
    "fmt"
    "Pact/corelib/client"
)

type Package struct {
	Name              string      `hcl:"name"`
	PackageIdentifier string      `hcl:"package_identifier"`
	Description       string      `hcl:"description"`
	InitRelease       InitRelease `hcl:"release,block"`
    Versioning        string      `hcl:"Versioning"` 
}

type InitRelease struct {
	Url     string `hcl:"url"`
	Version string `hcl:"version"`
	Hash    string `hcl:"hash"`
}


// Validate interprets a string s in the given base (0, 2 to 36) and bit size
// (0 to 64) and returns the corresponding value i.
func (p *Package) Validate(mode string) error {

    var rs client.RepositorySource

     switch mode {

    case "l":
        rs = client.NewLFilesystemSource("/somewhere")
    case "r":
        rs,_ = client.NewGithubSource("kasperjack/pact","main") //HE: skipping errors for now 


    default:
        panic(fmt.Sprintf("unknown validation mode %s", mode)) //REP: change this later 

    }
    
    repo := client.NewRepoClient(rs)

    // check package does not already exsit 
    // check versioning supported 
    // check init release uses samme versioning defined by pacakge 



	return nil
}
