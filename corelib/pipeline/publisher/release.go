package publisher


 import (
    "Pact/corelib/model"
	"Pact/corelib/client"
	"Pact/corelib/validate"

)


func NewRelease(r *model.Release) Pipeline {
	return &ReleasePipe{Model: r}
}

type ReleasePipe struct {
	Model *model.Release
	Client *client.RepoClient
}



func (r *ReleasePipe) Validate() error {
	return validate.Release(r.Model,r.Client)
}

func (r *ReleasePipe) Build() error {
	return nil
}