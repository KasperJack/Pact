package publisher

import (
	"Pact/corelib/client"
	"Pact/corelib/model"
	"Pact/corelib/pipeline"
	"Pact/corelib/validate"
)


func NewRelease(release  *model.Release, client *client.RepoClient ) pipeline.PublisherRelease {
	return &ReleaseContext{Model: release, Client: client}
}

type ReleaseContext  struct {
	Model *model.Release
	Client *client.RepoClient
}



func (rc *ReleaseContext) Validate() error {
	return validate.Release(rc.Model,rc.Client)
}

func (rc *ReleaseContext) Build() error {
	return nil
}