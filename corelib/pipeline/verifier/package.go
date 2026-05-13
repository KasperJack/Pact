package verifier


 import (
    "Pact/corelib/model"
	"Pact/corelib/client"
)

func NewPackage(p *model.Package) Pipeline {
	return &PackagePipe{Model: p}
}

type PackagePipe struct {
	Model *model.Package
	Client *client.RepoClient
}



func (p *PackagePipe) Validate() error {
	return nil
}

