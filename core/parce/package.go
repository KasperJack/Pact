package parce

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/kasperjack/pact/core"
	"regexp"
	"strings"


)

var validPackageNamePattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)




type lPackage struct {
	Range 	  hcl.Range

	Package       string   `hcl:"package"`
	Name          string   `hcl:"name"`
	Description   *string  `hcl:"description,optional"`
	Homepage      *string  `hcl:"homepage,optional"`
	License       *string  `hcl:"license,optional"`
	RawArchitectures []string `hcl:"architectures"`
	RawScopes        []string `hcl:"scopes"`

	Scopes        []core.Scope
	Architectures []core.Arch
}




func (p *lPackage) validate(body *hclsyntax.Body) hcl.Diagnostics {


	var allDiags hcl.Diagnostics



	rangeOf := func(field string) hcl.Range {

        return attrRangeOf(body, field)
    }




	allDiags = append(allDiags, checkField(p.Package, "package", rangeOf("package"))...)
    allDiags = append(allDiags, checkValidPackageName(p.Package, rangeOf("package"))...)


    allDiags = append(allDiags, checkField(p.Name, "name", rangeOf("name"))...)
    allDiags = append(allDiags, checkField(p.Description, "description", rangeOf("description"))...)
    allDiags = append(allDiags, checkField(p.Homepage, "homepage", rangeOf("homepage"))...)
    allDiags = append(allDiags, checkField(p.License, "license", rangeOf("license"))...)




	if len(p.RawArchitectures) == 0 {
		allDiags = append(allDiags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "architectures must not be empty",
			Detail:   "a package must declare at least one supported architecture",
			Subject:  rangeOf("architectures").Ptr(),
		})
	}


	if len(p.RawScopes) == 0 {
		allDiags = append(allDiags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "scopes must not be empty",
			Detail:   `a package must declare at least one of "user" or "system"`,
			Subject:  rangeOf("scopes").Ptr(),
		})
	}






	for _, raw := range p.RawArchitectures {

		_, err := core.ParseArch(raw)

		if err != nil {
			allDiags = append(allDiags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("invalid architecture %q", raw),
				Detail:   err.Error(),
				Subject:  rangeOf("architectures").Ptr(),
			})
			continue
		}

	}


	for _, raw := range p.RawScopes {
		_, err := core.ParseScope(raw)
		if err != nil {
			allDiags = append(allDiags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("invalid scope %q", raw),
				Detail:   err.Error(),
				Subject:  rangeOf("scopes").Ptr(),
			})
			continue
		}

	}



	allDiags = append(allDiags, checkDuplicateArchs(p.RawArchitectures, rangeOf("architectures"))...)
	allDiags = append(allDiags, checkDuplicateScopes(p.RawScopes, rangeOf("scopes"))...)



	return allDiags
}



func (p *lPackage) export() (*core.PackageInfo) {

	return &core.PackageInfo{
		PackageIdentifier: p.Package,
		Name: strings.TrimSpace(p.Name),
		Description: strings.TrimSpace(derefOr(p.Description,"")),
		Homepage: strings.TrimSpace(derefOr(p.Homepage,"")),
		License: strings.TrimSpace(derefOr(p.License,"")),

		Architectures: mustMapParsed(p.RawArchitectures, core.ParseArch),
        Scopes:        mustMapParsed(p.RawScopes, core.ParseScope),


	}



}
	















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



	pkg := &lPackage{}




	if d := gohcl.DecodeBody(f.Body, nil, pkg); d.HasErrors() {
		return nil, d
	}

	diags = pkg.validate(syntaxBody) //syntaxBody ? 

	



	if diags.HasErrors() {
		return nil, diags
	}


	return pkg.export(),diags
}







//pkg parce helpers



func checkValidPackageName(value string, rng hcl.Range) hcl.Diagnostics {

	if !validPackageNamePattern.MatchString(value) {
		return hcl.Diagnostics{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("invalid package name %q", value),
			Detail:   "package names must be lowercase letters, digits, and hyphens only, e.g. \"libtree\" or \"lib-tree\"",
			Subject:  rng.Ptr(),
		}}
	}
	return nil
}








func checkDuplicateArchs(archs []string, rng hcl.Range) hcl.Diagnostics {
	var diags hcl.Diagnostics

	seen := map[string]bool{}

	for _, a := range archs {
		if seen[a] {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("duplicate architecture %q", a),
				Detail:   "each architecture can be listed only once",
				Subject:  rng.Ptr(),
			})
			continue
		}
		seen[a] = true
	}

	return diags
}





func checkDuplicateScopes(scopes []string, rng hcl.Range) hcl.Diagnostics {
	var diags hcl.Diagnostics
	seen := map[string]bool{}

	for _, s := range scopes {
		if seen[s] {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("duplicate scope %q", s),
				Detail:   "each scope can be listed only once",
				Subject:  rng.Ptr(),
			})
			continue
		}
		seen[s] = true
	}

	return diags
}




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