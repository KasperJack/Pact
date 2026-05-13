package verifier


 import (
    "Pact/corelib/model"
	"Pact/corelib/client"
)

func NewRelease(r *model.Release) Pipeline {
	return &ReleasePipe{Model: r}
}

type ReleasePipe struct {
	Model *model.Release
	Client *client.RepoClient
}



func (r *ReleasePipe) Validate() error {
	
	return nil
}
