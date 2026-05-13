package publisher


 import (
    "Pact/corelib/model"
	"Pact/corelib/client"
	"Pact/corelib/validate"
)

func NewPackage(pacakge *model.Package, client *client.RepoClient) Pipeline {
	return &Package{Model: pacakge, Client: client}
}

type Package struct {
	Model *model.Package
	Client *client.RepoClient
}



func (p *Package) Validate() error {
	return validate.Package(p.Model,p.Client)
}


func (p *Package) Build() error {
	return nil
}