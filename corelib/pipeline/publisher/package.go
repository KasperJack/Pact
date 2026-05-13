package publisher


 import (
    "Pact/corelib/model"
	"Pact/corelib/client"
	"Pact/corelib/validate"
)

func NewPackage(p *model.Package, c *client.RepoClient) Pipeline {
	return &PackagePipe{Model: p, Client: c}
}

type PackagePipe struct {
	Model *model.Package
	Client *client.RepoClient
}



func (p *PackagePipe) Validate() error {
	return validate.Package(p.Model,p.Client)
}


func (p *PackagePipe) Build() error {
	return nil
}