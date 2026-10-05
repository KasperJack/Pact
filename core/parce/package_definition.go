package parce

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/kasperjack/pact/core"
	"strings"


)


// move the pkg name to a label 
type PackageDefinitionBlock struct {

	Package       string   `hcl:"package"`
	Name          string   `hcl:"name"`
	Description   *string  `hcl:"description,optional"`
	Homepage      *string  `hcl:"homepage,optional"`
	License       *string  `hcl:"license,optional"`
	RawArchitectures []string `hcl:"architectures"`
	RawScopes        []string `hcl:"scopes"`

}






func (p *PackageDefinitionBlock) export() (*core.PackageDefinition) {

	return &core.PackageDefinition{
		Identifier: p.Package,
		Name: strings.TrimSpace(p.Name),
		Description: strings.TrimSpace(derefOr(p.Description,"")),
		Homepage: strings.TrimSpace(derefOr(p.Homepage,"")),
		License: strings.TrimSpace(derefOr(p.License,"")),

		Architectures: mustMapParsed(p.RawArchitectures, core.ParseArch),
        Scopes:        mustMapParsed(p.RawScopes, core.ParseScope),


	}



}
	




func ParsePackageDefinition(src []byte, fileName string) (*core.PackageDefinition, hcl.Diagnostics) {

	parser := hclparse.NewParser()
	f, diags := parser.ParseHCL(src, fileName)

	if diags.HasErrors() {
		return nil, diags
	}


	pkg := &PackageDefinitionBlock{}

	// use the new genric method for blck decoding // it include handling of labels


	if d := gohcl.DecodeBody(f.Body, nil, pkg); d.HasErrors() {
		return nil, d
	}


	

	return pkg.export(),diags
}







//pkg parce helpers


func mustMapParsed[T any](raw []string, parse func(string) (T, error)) []T {
    out := make([]T, len(raw))
    for i, r := range raw {
        v, err := parse(r)
        if err != nil {
            panic(fmt.Sprintf("mapParsed: unexpected invalid value %q (was validate() called first?): %v", r, err))
        }
        out[i] = v
    }
    return out
}