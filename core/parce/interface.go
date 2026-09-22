package parce

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/kasperjack/pact/core"
	"strings"
)





func Interface(src []byte, acceptScope []core.Scope) (*core.Interface, hcl.Diagnostics) {
	parser := hclparse.NewParser()
	f, diags := parser.ParseHCL(src, "interface.hcl")
	if diags.HasErrors() {
		return nil, diags
	}

	syntaxBody, ok := f.Body.(*hclsyntax.Body)
	if !ok {
		return nil, hcl.Diagnostics{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "unexpected body type",
			Detail:   "could not parse interface as HCL native syntax",
		}}
	}

	if len(acceptScope) == 1 {
		return parseSingleScopeInterface(syntaxBody, acceptScope[0])
	}

	return parseMultiScopeInterface(syntaxBody, acceptScope)
}

// parseSingleScopeInterface handles the common case: the manifest only
// supports one scope, so there's nothing to disambiguate — option blocks
// live directly at the top level, no user{}/system{} wrapper needed.


func parseSingleScopeInterface(body *hclsyntax.Body, scope core.Scope) (*core.Interface, hcl.Diagnostics) {
	iface := &core.Interface{}
	var allDiags hcl.Diagnostics

	for _, block := range body.Blocks {
		switch block.Type {
		case "option":
			// reject scope-wrapper blocks here — if the manifest only
			// declared one scope, wrapping is meaningless overhead and
			// signals the publisher may be confused about which mode applies
			continue
		case "user", "system":
			allDiags = append(allDiags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("unexpected %q block", block.Type),
				Detail:   "this package only supports one scope — declare option blocks directly at the top level, without a user{}/system{} wrapper",
				Subject:  block.DefRange().Ptr(),
			})
			continue
		default:
			allDiags = append(allDiags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("unknown top-level block type %q", block.Type),
				Detail:   `only "option" blocks are allowed at the top level for a single-scope package`,
				Subject:  block.DefRange().Ptr(),
			})
		}
	}

	opts, d := parseOptionScope(body)
	allDiags = append(allDiags, d...)

	switch scope {
	case core.ScopeUser:
		iface.User = opts
	case core.ScopeSystem:
		iface.System = opts
	}

	if allDiags.HasErrors() {
		return nil, allDiags
	}
	return iface, allDiags
}

// parseMultiScopeInterface handles the dual-scope case — options may need
// to differ per scope, so the user{}/system{} wrapper is required here to
// disambiguate which set applies to which resolved scope.
func parseMultiScopeInterface(body *hclsyntax.Body, acceptScope []core.Scope) (*core.Interface, hcl.Diagnostics) {
	iface := &core.Interface{}
	var allDiags hcl.Diagnostics
	seen := map[string]bool{}

	for _, block := range body.Blocks {
		switch block.Type {
		case "user":
			seen["user"] = true
			opts, d := parseOptionScope(block.Body)
			allDiags = append(allDiags, d...)
			iface.User = opts

		case "system":
			seen["system"] = true
			opts, d := parseOptionScope(block.Body)
			allDiags = append(allDiags, d...)
			iface.System = opts

		case "option":
			allDiags = append(allDiags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "unexpected top-level option block",
				Detail:   "this package supports multiple scopes — wrap option blocks in user{} and/or system{}",
				Subject:  block.DefRange().Ptr(),
			})

		default:
			allDiags = append(allDiags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("unknown top-level block type %q", block.Type),
				Detail:   `only "user" and "system" blocks are allowed at the top level`,
				Subject:  block.DefRange().Ptr(),
			})
		}
	}

	// require a wrapper block for every scope the manifest actually declared —
	// same "must match or be explained" discipline as scope itself
	for _, s := range acceptScope {
		name := scopeBlockName(s)
		if !seen[name] {
			allDiags = append(allDiags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("missing %q block", name),
				Detail:   fmt.Sprintf("package declares %q as a supported scope, but interface.hcl has no %s{} block", name, name),
				Subject:  body.SrcRange.Ptr(),
			})
		}
	}

	if allDiags.HasErrors() {
		return nil, allDiags
	}
	return iface, allDiags
}

func scopeBlockName(s core.Scope) string {
	switch s {
	case core.ScopeUser:
		return "user"
	case core.ScopeSystem:
		return "system"
	default:
		return "unknown"
	}
}
















func parseOptionScope(body *hclsyntax.Body) ([]core.Option, hcl.Diagnostics) {
	var opts []core.Option
	var diags hcl.Diagnostics

	for _, block := range body.Blocks {
		if block.Type != "option" {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("unknown block type %q", block.Type),
				Detail:   `only "option" blocks are allowed here`,
				Subject:  block.DefRange().Ptr(),
			})
			continue
		}

		opt, d := parseOption(block)
		diags = append(diags, d...)
		if d.HasErrors() {
			continue
		}
		opts = append(opts, opt)
	}

	return opts, diags
}

func parseOption(block *hclsyntax.Block) (core.Option, hcl.Diagnostics) {
	id, diags := optionLabel(block) // same 0-or-1-label + validIDPattern rule as blockLabel
	if diags.HasErrors() {
		return core.Option{}, diags
	}

	var attrs struct {
		Default     bool     `hcl:"default"`
		Label       *string  `hcl:"label,optional"`
		Description *string  `hcl:"description,optional"`
		Binding     []string `hcl:"binding"`
	}
	if d := gohcl.DecodeBody(block.Body, nil, &attrs); d.HasErrors() {
		return core.Option{}, d
	}

	diags = append(diags, checkOptional(attrs.Label, "label", attrRangeOf(block, "label"))...)
	diags = append(diags, checkOptional(attrs.Description, "description", attrRangeOf(block, "description"))...)

	if len(attrs.Binding) == 0 {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("option %q must have at least one binding", id),
			Detail:   "binding must not be an empty list",
			Subject:  attrRangeOf(block, "binding").Ptr(),
		})
	}

	if diags.HasErrors() {
		return core.Option{}, diags
	}

	return core.Option{
		ID:          id,
		Default:     attrs.Default,
		Label:       strings.TrimSpace(derefOr(attrs.Label, "")),
		Description: strings.TrimSpace(derefOr(attrs.Description, "")),
		Binding:     attrs.Binding,
	}, diags
}

// optionLabel mirrors blockLabel's 0-or-1-label rule, but requires exactly 1
// (options are never unlabeled), reusing validIDPattern.
func optionLabel(block *hclsyntax.Block) (string, hcl.Diagnostics) {
	if len(block.Labels) != 1 {
		return "", hcl.Diagnostics{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("option requires exactly 1 label, got %d", len(block.Labels)),
			Detail:   `e.g. option "desktop_shortcut" { ... }`,
			Subject:  block.DefRange().Ptr(),
		}}
	}

	label := block.Labels[0]
	if strings.TrimSpace(label) == "" {
		return "", hcl.Diagnostics{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "option has an empty id",
			Subject:  block.DefRange().Ptr(),
		}}
	}
	if !validIDPattern.MatchString(label) {
		return "", hcl.Diagnostics{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("invalid id %q", label),
			Detail:   "ids may only contain letters, digits, underscore, and hyphen — no whitespace",
			Subject:  block.DefRange().Ptr(),
		}}
	}
	return label, nil
}