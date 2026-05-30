package verifier

 import (
    "Pact/corelib/model"
	"Pact/corelib/pipeline"
	"Pact/corelib/client"
	"Pact/corelib/validate"

)


func NewPackage(pacakge *model.Package, client *client.RepoClient) pipeline.VerifierPackage {
	return &PackageContext{Model: pacakge, Client: client}
}

type PackageContext struct {
	Model *model.Package
	Client *client.RepoClient
}



func (p *PackageContext) Validate() error {
	return validate.Package(p.Model,p.Client)
}

