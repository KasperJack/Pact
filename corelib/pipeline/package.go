package pipeline


 import (

    "Pact/corelib/model"
)

func NewPackage(p *model.Package) Pipeline {
	return &PackagePipe{Model: p}
}

type PackagePipe struct {
	Model *model.Package
}



func (p *PackagePipe) Validate() error {
	return nil
}

func (p *PackagePipe) Build() error {
	return nil
}