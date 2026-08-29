package parce

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/kasperjack/pact/core"
)



func PackageInfo(src []byte) (*core.PackageInfo, hcl.Diagnostics) {

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

	var attrs struct {
		Package       string   `hcl:"package"`
		Name          string   `hcl:"name"`
		Description   *string  `hcl:"description,optional"`
		Homepage      *string  `hcl:"homepage,optional"`
		License       *string  `hcl:"license,optional"`
		Architectures []string `hcl:"architectures"`
		Scopes        []string `hcl:"scopes"`
	}
	if d := gohcl.DecodeBody(f.Body, nil, &attrs); d.HasErrors() {
		return nil, d
	}

	var allDiags hcl.Diagnostics

	allDiags = append(allDiags, checkRequired(attrs.Package, "package", attrRangeOfBody(syntaxBody, "package"))...)
	allDiags = append(allDiags, checkRequired(attrs.Name, "name", attrRangeOfBody(syntaxBody, "name"))...)
	allDiags = append(allDiags, checkOptional(attrs.Description, "description", attrRangeOfBody(syntaxBody, "description"))...)
	allDiags = append(allDiags, checkOptional(attrs.Homepage, "homepage", attrRangeOfBody(syntaxBody, "homepage"))...)
	allDiags = append(allDiags, checkOptional(attrs.License, "license", attrRangeOfBody(syntaxBody, "license"))...)

	if len(attrs.Architectures) == 0 {
		allDiags = append(allDiags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "architectures must not be empty",
			Detail:   "a package must declare at least one supported architecture",
			Subject:  attrRangeOfBody(syntaxBody, "architectures").Ptr(),
		})
	}
	if len(attrs.Scopes) == 0 {
		allDiags = append(allDiags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "scopes must not be empty",
			Detail:   `a package must declare at least one of "user" or "system"`,
			Subject:  attrRangeOfBody(syntaxBody, "scopes").Ptr(),
		})
	}

	archs := make([]core.Arch, 0, len(attrs.Architectures))
	for _, raw := range attrs.Architectures {
		a, err := core.ParseArch(raw)
		if err != nil {
			allDiags = append(allDiags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("invalid architecture %q", raw),
				Detail:   err.Error(),
				Subject:  attrRangeOfBody(syntaxBody, "architectures").Ptr(),
			})
			continue
		}
		archs = append(archs, a)
	}

	scopes := make([]core.Scope, 0, len(attrs.Scopes))
	for _, raw := range attrs.Scopes {
		s, err := core.ParseScope(raw)
		if err != nil {
			allDiags = append(allDiags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("invalid scope %q", raw),
				Detail:   err.Error(),
				Subject:  attrRangeOfBody(syntaxBody, "scopes").Ptr(),
			})
			continue
		}
		scopes = append(scopes, s)
	}

	if allDiags.HasErrors() {
		return nil, allDiags
	}

	return &core.PackageInfo{
		Package:       strings.TrimSpace(attrs.Package),
		Name:          strings.TrimSpace(attrs.Name),
		Description:   strings.TrimSpace(derefOr(attrs.Description, "")),
		Homepage:      strings.TrimSpace(derefOr(attrs.Homepage, "")),
		License:       strings.TrimSpace(derefOr(attrs.License, "")),
		Architectures: archs,
		Scopes:        scopes,
	}, allDiags
}

// attrRangeOfBody mirrors attrRangeOf but works on a *hclsyntax.Body directly
// (top-level file attributes) instead of a block's body.
func attrRangeOfBody(body *hclsyntax.Body, name string) hcl.Range {
	if attr, ok := body.Attributes[name]; ok {
		return attr.Expr.Range()
	}
	return body.SrcRange
}