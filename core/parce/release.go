package parce

import (
	
	"github.com/kasperjack/pact/core"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"

)



type release struct {


    Package         string `hcl:"package"`
    UpstreamVersion string `hcl:"upstream_version"` 
    Revision        int     `hcl:"revision"` 
    URL             string  `hcl:"url"` 
    SHA256          string `hcl:"sha256"` 
    SizeMB          int  `hcl:"size_mb"` 

    ArchitectureRaw string `hcl:"architecture"`
    
    Architecture core.Arch 
}





func (*release) export () *core.Release {

	return nil
}


func (*release) validate(body *hclsyntax.Body) hcl.Diagnostics {
	return nil
}








func Release(src []byte) (*core.Release,error) {


	parser := hclparse.NewParser()
	f, diags := parser.ParseHCL(src, "package.hcl")

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

	diags = release.validate(syntaxBody) //syntaxBody ? 

	



	if diags.HasErrors() {
		return nil, diags
	}


	return release.export(),diags
}




