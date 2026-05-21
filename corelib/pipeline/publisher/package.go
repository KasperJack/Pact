package publisher


 import (
    "Pact/corelib/model"
	"Pact/corelib/client"
	"Pact/corelib/validate"
	"Pact/corelib/build"
	"Pact/corelib/pipeline"
	"Pact/corelib/index"
)

func NewPackage(pacakge *model.Package, client *client.RepoClient) pipeline.PublisherPackage {
	return &PackageContext{Model: pacakge, Client: client}
}

type PackageContext struct {  // Context ? change name 
	Model *model.Package
	Client *client.RepoClient
}


// Package
func (p *PackageContext) Validate() error {
	
	return validate.Package(p.Model,p.Client)
}


func (p *PackageContext) Build() error {
	return build.Package(p.Model,p.Client.GetRoot())
}

func (p *PackageContext) RebuildIndex() error {
	return index.Rebuild()
}