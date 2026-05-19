package validate

import (
    "fmt"
    "Pact/corelib/client"
    "Pact/corelib/model"
)




func Package(p *model.Package, c *client.RepoClient) error {

    ok := c.PackageExists(p.PackageIdentifier)

    if ok {
        // error packag already exists
        return fmt.Errorf("Package already exists: %s", p.PackageIdentifier)
    }

    switch p.Versioning {

    case "semver":
        if !isValidSemver(p.InitRelease.Version) {
            return fmt.Errorf("invalid semver version: %s", p.InitRelease.Version) }
    case "date":
        //
    case "custom":
        //

    default:
        return fmt.Errorf("invalid Versioning: %s", p.Versioning)
    }
    

    fmt.Println("package validated ok satus")
    //fmt.Println(c.PackageExists("windirstat"))
	return nil
}

