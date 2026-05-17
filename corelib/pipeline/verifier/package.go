package verifier

 import (
    "Pact/corelib/model"
	"Pact/corelib/pipeline"
	"Pact/corelib/client"
	"Pact/corelib/validate"

)


func NewPackage(pacakge *model.Package, client *client.RepoClient) pipeline.VerifierPackage {
	return &Package{Model: pacakge, Client: client}
}

type Package struct {
	Model *model.Package
	Client *client.RepoClient
}



func (p *Package) Validate() error {
	return validate.Package(p.Model,p.Client)
}

