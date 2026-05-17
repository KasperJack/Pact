package validate

import (
    "fmt"
    "Pact/corelib/client"
    "Pact/corelib/model"
)




func Package(p *model.Package, c *client.RepoClient) error {
    /*
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
    */

    fmt.Println("validating pcakge")
    fmt.Println(c.PackageExists("windirstat"))
	return nil
}

