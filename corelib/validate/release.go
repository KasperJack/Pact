package validate

import (
    "fmt"
    "Pact/corelib/client"
    "Pact/corelib/model"
)





func Release(p *model.Release, c *client.RepoClient) error {

    fmt.Println("validating releace")
	return nil
}
