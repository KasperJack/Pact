package publisher


 import (
    "Pact/corelib/model"
	"Pact/corelib/client"
	"Pact/corelib/validate"

)


func NewRelease(release  *model.Release, client *client.RepoClient ) Pipeline {
	return &Release{Model: release, Client: client}
}

type Release struct {
	Model *model.Release
	Client *client.RepoClient
}



func (r *Release) Validate() error {
	return validate.Release(r.Model,r.Client)
}

func (r *Release) Build() error {
	return nil
}