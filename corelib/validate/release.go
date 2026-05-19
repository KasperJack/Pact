package validate

import (
    //"fmt"
    "Pact/corelib/client"
    "Pact/corelib/model"
)





func Release(r *model.Release, c *client.RepoClient) error {

    ok := c.PackageExists(r.PackageIdentifier)
    if ok == false {
        // error the target pcakge does not exist 
    }
    // get versining schema of the pcakge 
    // make sure r.Version is a valid secema
    // make sure r.Version not existing pcakge versions 


	return nil
}
