package publisher

import (
	"Pact/corelib/build"
	"Pact/corelib/client"
	"Pact/corelib/model"
	"Pact/corelib/pipeline"
	"Pact/corelib/validate"
	"fmt"
	//"Pact/corelib/index"
	"Pact/corelib/luautil"
)


func NewPackage(luaData []byte, client *client.RepoClient) (pipeline.PublisherPackage, error) {

	var pacakge model.Package
    ctx := luautil.NewPackageEvalContext(&pacakge)
    defer ctx.Close()

    if err := ctx.Eval(luaData); err != nil {
        return nil, err
    }

	
	fmt.Println(pacakge.PackageIdentifier)
	fmt.Println(pacakge.Name)
	fmt.Println(pacakge.Versioning)

	fmt.Println(pacakge.Description)
	fmt.Println(pacakge.Homepage)
	fmt.Println(pacakge.License)
	
	err := ctx.RunInstall()

	if err != nil {
		return nil,err
	}

    return &PackageContext{Model: &pacakge, Client: client}, nil
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
	return nil
}