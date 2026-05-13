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



func (p *Package) Validate(mode string) error {

    switch p.Versioning {
    case "semver":
       //p.InitRelease.Version == "1.2.3"
    case "date":
        //
    case "custom":
        //

    default:
        return fmt.Errorf("invalid Versioning: %s", p.Versioning)
    }





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

    ok := repo.PackageExists(p.PackageIdentifier)
    if ok {
        return fmt.Errorf(fmt.Sprintf("package %s already exists",p.PackageIdentifier))
    }


    // check package does not already exsit 
    // check versioning supported 
    // check init release uses samme versioning defined by pacakge 



	return nil
}


de