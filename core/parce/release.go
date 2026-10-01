package parce

import (
	
	"github.com/kasperjack/pact/core"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"fmt"

)



type release struct {


    Package         string `hcl:"package"`

    UpstreamVersion string `hcl:"upstream_version" validate:"nonempty"` 

    Revision        int     `hcl:"revision"` 

    URL             string  `hcl:"url" validate:"nonempty"` 
    SHA256          string `hcl:"sha256" validate:"nonempty"` 
    SizeMB          int  `hcl:"size_mb"` 

	
    ArchitectureRaw string `hcl:"architecture"`
    Architecture core.Arch

}







func (r *release) validate(body *hclsyntax.Body) hcl.Diagnostics {

	var allDiags hcl.Diagnostics



	rangeOf := func(field string) hcl.Range {

        return attrRangeOf(body, field)
    }


	allDiags = append(allDiags, checkValidPackageIdentifier(r.Package, rangeOf("package"))...)
	allDiags = append(allDiags, validateStringFields(r, body)...)
	



	_, err := core.ParseArch(r.ArchitectureRaw)

		if err != nil {
			allDiags = append(allDiags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("invalid architecture %q", r.ArchitectureRaw),
				Detail:   err.Error(),
				Subject:  rangeOf("architecture").Ptr(),
			})
		}





	return allDiags
}





func (r *release) export() *core.Release {
	arch, _ := core.ParseArch(r.ArchitectureRaw)

	return &core.Release{
		PackageIdentifier: r.Package,
		UpstreamVersion:   r.UpstreamVersion,
		Revision:          r.Revision,
		URL:               r.URL,
		SHA256:             r.SHA256,
		SizeMB:             r.SizeMB,
		Architecture:      arch,
	}
}







func Release(src []byte) (*core.Release, hcl.Diagnostics) {


	parser := hclparse.NewParser()
	f, diags := parser.ParseHCL(src, "release.hcl")

	if diags.HasErrors() {
		return nil, diags
	}

	syntaxBody, ok := f.Body.(*hclsyntax.Body)
	if !ok {
		return nil, hcl.Diagnostics{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "unexpected body type",
			Detail:   "could not parse package.hcl as HCL native syntax",
		}}
	}



	release := &release{}




	if d := gohcl.DecodeBody(f.Body, nil, release); d.HasErrors() {
		return nil, d
	}




	fmt.Println("calling vlidate")
	diags = release.validate(syntaxBody) 
	fmt.Println("vlidate done")

	



	if diags.HasErrors() {
		return nil, diags
	}




	return release.export(),diags
}




