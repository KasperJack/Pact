package verifier

import (
	"Pact/corelib/client"
	"Pact/corelib/model"
	"Pact/corelib/pipeline"
	"Pact/corelib/validate"
)


func NewRelease(release  *model.Release, client *client.RepoClient ) pipeline.VerifierRelease {
	return &ReleaseContext{Model: release, Client: client}
}

type ReleaseContext struct {
	Model *model.Release
	Client *client.RepoClient
}



func (r *ReleaseContext) Validate() error {
	return validate.Release(r.Model,r.Client)
}