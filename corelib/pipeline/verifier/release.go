package verifier


 import (
	"fmt"
    "Pact/corelib/model"
)

func NewRelease(r *model.Release) Pipeline {
	return &ReleasePipe{Model: r}
}

type ReleasePipe struct {
	Model *model.Release
}



func (r *ReleasePipe) Validate() error {
	fmt.Printf("validationg relase %s",r.Model.Version)
	return nil
}
