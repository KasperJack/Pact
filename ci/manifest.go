package main

import (
	"fmt"
	"strings"
	"github.com/hashicorp/hcl/v2"
	"github.com/kasperjack/pact/core"
	"regexp"
)



var (
	invalidPathChars  = regexp.MustCompile(`[<>"|?*]`)
)


type ValidatedManifest struct {
	User   *core.ResolvedScope
	System *core.ResolvedScope
}




func ValidateManifest(m *core.ResolvedManifest) (*ValidatedManifest, hcl.Diagnostics) {


	var diags hcl.Diagnostics

	if m.User != nil {
		diags = append(diags, validateScope(m.User)...)
	}
	if m.System != nil {
		diags = append(diags, validateScope(m.System)...)
	}

	if diags.HasErrors() {
		return nil, diags
	}
	return &ValidatedManifest{User: m.User, System: m.System}, nil
}




func validateScope(scope *core.ResolvedScope) hcl.Diagnostics {
	var diags hcl.Diagnostics

	shortcutIDs := map[string]hcl.Range{}
	shortcutExes := map[string]hcl.Range{}
	//shortcutDisplayName := map[string]hcl.Range{}

	commandIDs := map[string]hcl.Range{}
	commandExes := map[string]hcl.Range{}

	addPathIDs := map[string]hcl.Range{}
	addPathDirs := map[string]hcl.Range{}



	for _, b := range scope.Blocks {



		switch v := b.(type) {
		case core.Shortcut:
			diags = append(diags, validateShortcut(v, shortcutIDs, shortcutExes)...)
		case core.Command:
			diags = append(diags, validateCommand(v, commandIDs,commandExes)...)
		case core.AddPath:
			diags = append(diags, validateAddPath(v, addPathIDs,addPathDirs)...)
		}
	}

	return diags
}





func validateShortcut(s core.Shortcut, ids, exes map[string]hcl.Range) hcl.Diagnostics {
	var diags hcl.Diagnostics

	diags = append(diags, checkDup(ids, s.ID, s.Range, "shortcut id")...)
	diags = append(diags, checkDup(exes, s.Exe, s.Range, "shortcut exe")...)
	diags = append(diags, checkValidPath(s.Exe, "exe", s.Range)...)
	diags = append(diags, checkValidPath(s.Icon, "icon", s.Range)...)

	return diags
}

func validateCommand(c core.Command, ids, exes map[string]hcl.Range) hcl.Diagnostics {
	var diags hcl.Diagnostics

	diags = append(diags, checkDup(ids, c.ID, c.Range, "command id")...)
	diags = append(diags, checkDup(exes, c.Exe, c.Range, "command exe")...)
	diags = append(diags, checkValidPath(c.Exe, "exe", c.Range)...)

	return diags
}

func validateAddPath(a core.AddPath, ids, dirs map[string]hcl.Range) hcl.Diagnostics {
	var diags hcl.Diagnostics

	diags = append(diags, checkDup(ids, a.ID, a.Range, "add_path id")...)
	diags = append(diags, checkDup(dirs, a.Dir, a.Range, "add_path dir")...)
	diags = append(diags, checkValidPath(a.Dir, "dir", a.Range)...)

	return diags
}









func checkDup(seen map[string]hcl.Range, value string, rng hcl.Range, what string) hcl.Diagnostics {
	if value == "" {
		return nil
	}
	if prev, exists := seen[value]; exists {
		return hcl.Diagnostics{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("duplicate %s", what),
			Detail:   fmt.Sprintf("%q was already used at %s; each %s must be unique within its scope", value, prev.String(), what),
			Subject:  rng.Ptr(),
		}}
	}
	seen[value] = rng
	return nil
}




func checkValidPath(value, field string, rng hcl.Range) hcl.Diagnostics {
	if value == "" {
		return nil // required-ness is checked separately; nothing to validate here
	}
	if invalidPathChars.MatchString(value) {
		return hcl.Diagnostics{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("%s is not a valid path", field),
			Detail:   fmt.Sprintf("%q contains characters not allowed in a Windows path: < > \" | ? *", value),
			Subject:  rng.Ptr(),
		}}
	}
	if strings.Contains(value, "..") {
		return hcl.Diagnostics{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("%s must not contain path traversal", field),
			Detail:   fmt.Sprintf("%q contains \"..\", which is not allowed", value),
			Subject:  rng.Ptr(),
		}}
	}
	return nil
}

