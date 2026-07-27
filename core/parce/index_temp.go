package parce

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"

)



func GetType(pkgData []byte) (string, error) {
	attr, err := getAttr(pkgData, "type")
	if err != nil {
		return "", err
	}

	var value string
	diags := gohcl.DecodeExpression(attr.Expr, nil, &value)
	if diags.HasErrors() {
		return "", diags
	}
	return value, nil
}




func GetLatest(pkgData []byte) (string, error) {
	attr, err := getAttr(pkgData, "latest")
	if err != nil {
		return "", err
	}

	var value string
	diags := gohcl.DecodeExpression(attr.Expr, nil, &value)
	if diags.HasErrors() {
		return "", diags
	}
	return value, nil
}

func GetVersions(pkgData []byte) ([]string, error) {
	attr, err := getAttr(pkgData, "versions")
	if err != nil {
		return nil, err
	}

	var value []string
	diags := gohcl.DecodeExpression(attr.Expr, nil, &value)
	if diags.HasErrors() {
		return nil, diags
	}
	return value, nil
}

// getAttr is a shared helper parces the file and pulls one named attribute
func getAttr(pkgData []byte, name string) (*hcl.Attribute, error) {
	parser := hclparse.NewParser()
	f, diags := parser.ParseHCL(pkgData, "[package].hcl")
	if diags.HasErrors() {
		return nil, diags
	}

	schema := &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{{Name: name}},
	}
	content, _, diags := f.Body.PartialContent(schema)
	if diags.HasErrors() {
		return nil, diags
	}

	attr, ok := content.Attributes[name]
	if !ok {
		return nil, fmt.Errorf("%s: attribute not found", name)
	}
	return attr, nil
}









func GetSourceBlockNames(src []byte) ([]string, error) {
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL(src, "ass")
	if diags.HasErrors() {
		return nil, fmt.Errorf("parsing %s: %w", "ass", diags)
	}

	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		return nil, fmt.Errorf("unexpected body type for %s", "ass")
	}

	var sourceBody *hclsyntax.Body
	for _, block := range body.Blocks {
		if block.Type == "source" {
			sourceBody = block.Body
			break
		}
	}
	if sourceBody == nil {
		return nil, fmt.Errorf("no source block found in %s", "ass")
	}

	names := make([]string, 0, len(sourceBody.Blocks))
	for _, archBlock := range sourceBody.Blocks {
		names = append(names, archBlock.Type)
	}

	return names, nil
}