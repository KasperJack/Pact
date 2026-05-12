package stage

import (
    "fmt"

)

type Package struct {
	Name              string      `hcl:"name"`
	PackageIdentifier string      `hcl:"package_identifier"`
	Description       string      `hcl:"description"`
	InitRelease       InitRelease `hcl:"release,block"`
}

type InitRelease struct {
	Url     string `hcl:"url"`
	Version string `hcl:"version"`
	Hash    string `hcl:"hash"`
}

func (p *Package) Validate(mode string) error {
	fmt.Println(p.Name)
    fmt.Println(p.PackageIdentifier)
    fmt.Println(p.Description)
    fmt.Println(p.InitRelease.Url)
    fmt.Println(p.InitRelease.Version)
    fmt.Println(p.InitRelease.Hash)
	return nil
}
