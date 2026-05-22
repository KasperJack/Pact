package validate

import (
    "fmt"
    "Pact/corelib/client"
    "Pact/corelib/model"
)




func Package(p *model.Package, c *client.RepoClient) error {

    ok := c.PackageExists(p.PackageIdentifier)
    if ok {
        return fmt.Errorf("package %q already exists", p.PackageIdentifier)
    }

    ok = isValidVersioningSchema(p.Versioning)
    if !ok {
        return fmt.Errorf("invalid versioning schema %q", p.Versioning)
    }

    ok = isValidVersion(p.Versioning, p.Versioning) //??????
    if !ok {
        return fmt.Errorf("invalid version %q for versioning schema %q", p.Versioning, p.Versioning) //??????
    }

    fmt.Println("package validated ok satus")
    //fmt.Println(c.PackageExists("windirstat"))
	return nil
}

