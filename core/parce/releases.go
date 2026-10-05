package parce





import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/kasperjack/pact/core"


)




// use the new genric method for blck decoding 
func Release(src []byte, fileName string) (*core.Release, hcl.Diagnostics) {

	r := &core.Release{}

	
	parser := hclparse.NewParser()
	f, diags := parser.ParseHCL(src, fileName)
	if diags.HasErrors() {
		return nil, diags
	}



	syntaxBody, ok := f.Body.(*hclsyntax.Body)

	if !ok {
		return nil, hcl.Diagnostics{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "unexpected body type",
			Detail:   "could not parse manifest.hcl as HCL native syntax",
		}}
	}

	
	relaseBlock := syntaxBody.Blocks[0]


	

	for _, block := range relaseBlock.Body.Blocks {

		switch block.Type {

		case "package":
			pkg, d := parsePackage(block)
			diags = append(diags, d...)
			r.Package = pkg


		case "source":
			src, d := parseSource(block)
			diags = append(diags, d...)
			r.Source = src

		case "state":
			Resources, d := parseState(block)
			diags = append(diags, d...)
			r.Resources = Resources


		default:
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("unknown block type %q", block.Type),
				Detail:   `expected "package", "source", or "state"`,
				Subject:  block.DefRange().Ptr(),
			})
		}


	}



	return r, diags

}