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

func Interface(src []byte) (*core.Interface, hcl.Diagnostics) {
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

	iface := &core.Interface{}
	var allDiags hcl.Diagnostics
	found := false

	for _, block := range syntaxBody.Blocks {
		switch block.Type {
		case "user":
			found = true
			opts, d := parseOptionScope(block.Body)
			allDiags = append(allDiags, d...)
			iface.User = opts

		case "system":
			found = true
			opts, d := parseOptionScope(block.Body)
			allDiags = append(allDiags, d...)
			iface.System = opts

		default:
			allDiags = append(allDiags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("unknown top-level block type %q", block.Type),
				Detail:   `only "user" and "system" blocks are allowed at the top level`,
				Subject:  block.DefRange().Ptr(),
			})
		}
	}

	if !found {
		allDiags = append(allDiags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "interface declares no user or system block",
			Detail:   "an interface file must declare at least one of user{} or system{}",
			Subject:  syntaxBody.SrcRange.Ptr(),
		})
	}

	if allDiags.HasErrors() {
		return nil, allDiags
	}
	return iface, allDiags
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
	c := core.Common{ID: id, Range: block.DefRange()}

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
		Common:      c,
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