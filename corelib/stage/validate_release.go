package stage

import (
	"fmt"


)

type Release struct {
    PackageIdentifier string `hcl:"package_identifier"`
	Url    string `hcl:"url"`
	Version string `hcl:"version"`
    Hash string `hcl:"hash"`
}

func (r *Release) Validate () error {
    fmt.Println("validating release")
    return nil
}  
