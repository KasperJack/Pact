package parce

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/kasperjack/pact/core"

	"strings"
	"regexp"

	//"/github.com/zclconf/go-cty/cty"
)



var (
	validIDPattern    = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	//invalidPathChars  = regexp.MustCompile(`[<>"|?*]`)
)


func Manifest(src []byte) (*core.ResolvedManifest, hcl.Diagnostics) {
	parser := hclparse.NewParser()
	f, diags := parser.ParseHCL(src, "manifest.hcl")
	if diags.HasErrors() {
		return nil, diags
	}

	syntaxBody, ok := f.Body.(*hclsyntax.Body)
	if !ok {
		return nil, hcl.Diagnostics{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "unexpected body type",
			Detail:   fmt.Sprintf("could not parse %s as HCL native syntax", "manifest.hcl"),
		}}
	}

	m := &core.ResolvedManifest{}
	found := false

	for _, block := range syntaxBody.Blocks {
		switch block.Type {

		case "user":
			found = true
			if m.User != nil {
				diags = append(diags, dupTopLevelErr("user", block.DefRange()))
				continue
			}
			scope, d := parseScope(block.Body)
			diags = append(diags, d...)
			m.User = scope

		case "system":
			found = true
			if m.System != nil {
				diags = append(diags, dupTopLevelErr("system", block.DefRange()))
				continue
			}
			scope, d := parseScope(block.Body)
			diags = append(diags, d...)
			m.System = scope

		default:
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("unknown top-level block type %q", block.Type),
				Detail:   `only "user" and "system" blocks are allowed at the top level`,
				Subject:  block.DefRange().Ptr(),
			})
		}
	}

	if !found {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "manifest declares no user or system block",
			Detail:   "a manifest must declare at least one of user{} or system{}",
			Subject:  syntaxBody.SrcRange.Ptr(),
		})
	}

	if diags.HasErrors() {
		return nil, diags
	}
	return m, diags
}



func dupTopLevelErr(blockType string, rng hcl.Range) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  fmt.Sprintf("duplicate %s block", blockType),
		Detail:   fmt.Sprintf("a manifest may declare at most one %s block", blockType),
		Subject:  rng.Ptr(),
	}
}

// ---------- scope-level parse ----------

func parseScope(body *hclsyntax.Body) (*core.ResolvedScope, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	installPath, d := decodeInstallPath(body)
	diags = append(diags, d...)

	scope := &core.ResolvedScope{InstallPath: installPath}

	for _, block := range body.Blocks {
		b, d := parseBlock(block)
		diags = append(diags, d...)
		if b != nil {
			scope.Blocks = append(scope.Blocks, b) // single append point = order preserved
		}
	}

	return scope, diags
}

func decodeInstallPath(body *hclsyntax.Body) (string, hcl.Diagnostics) {
	attr, ok := body.Attributes["install_path"]
	if !ok {
		return "", hcl.Diagnostics{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "missing required attribute install_path",
			Subject:  body.SrcRange.Ptr(),
		}}
	}

	val, diags := attr.Expr.Value(nil)

	if diags.HasErrors() {
		return "", diags
	}
	if val.Type().FriendlyName() != "string" {
		return "", hcl.Diagnostics{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "install_path must be a string",
			Subject:  attr.Expr.Range().Ptr(),
		}}
	}


	installPath := val.AsString()

	

	checkDiags := checkRequired(installPath, "install_path", attr.Expr.Range())
	if checkDiags.HasErrors() {
		return "", checkDiags
	}



	return strings.TrimSpace(installPath), nil
}

// ---------- block dispatch ----------

func parseBlock(block *hclsyntax.Block) (core.ResolvedBlock, hcl.Diagnostics) {

	id, diags := blockLabel(block)

	if diags.HasErrors() {
		return nil, diags
	}

	c := core.Common{ID: id, Range: block.DefRange()}

	switch block.Type {
	case "shortcut":
		s, d := parseShortcut(block, c)
		return s, d
	case "command":
		cmd, d := parseCommand(block, c)
		return cmd, d
	case "add_path":
		a, d := parseAddPath(block, c)
		return a, d
	default:
		return nil, hcl.Diagnostics{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("unknown block type %q", block.Type),
			Detail:   `expected "shortcut", "command", or "add_path"`,
			Subject:  block.DefRange().Ptr(),
		}}
	}
}

// blockLabel enforces the 0-or-1-label rule and rejects whitespace-only labels.
func blockLabel(block *hclsyntax.Block) (string, hcl.Diagnostics) {
	switch len(block.Labels) {
	case 0:
		return "", nil
	case 1:
		label := block.Labels[0]
		if strings.TrimSpace(label) == "" {
			return "", hcl.Diagnostics{&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("block %q has an empty label", block.Type),
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


	default:
		return "", hcl.Diagnostics{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("block %q takes 0 or 1 labels, got %d", block.Type, len(block.Labels)),
			Subject:  block.DefRange().Ptr(),
		}}
	}
}

// ---------- per-type parse: HCL decode shapes live only here ----------




func parseShortcut(block *hclsyntax.Block, c core.Common) (core.Shortcut, hcl.Diagnostics) {
	var attrs struct {
		DisplayName *string `hcl:"display_name,optional"`
		Exe         string  `hcl:"exe"`
		Icon        *string `hcl:"icon,optional"`
		Args        *string `hcl:"args,optional"`
	}


	if diags := gohcl.DecodeBody(block.Body, nil, &attrs); diags.HasErrors() {
		return core.Shortcut{}, diags
	}


	var diags hcl.Diagnostics

	diags = append(diags, checkRequired(attrs.Exe, "exe", block.Body.Attributes["exe"].Expr.Range())...)
	diags = append(diags, checkOptional(attrs.DisplayName, "display_name", attrRangeOf(block, "display_name"))...)
	diags = append(diags, checkOptional(attrs.Icon, "icon", attrRangeOf(block, "icon"))...)
	diags = append(diags, checkOptional(attrs.Args, "args", attrRangeOf(block, "args"))...)


	if diags.HasErrors() {
		return core.Shortcut{}, diags
	}

	return core.Shortcut{
		Common:      c,
		Exe:         strings.TrimSpace(attrs.Exe),
		DisplayName: strings.TrimSpace(derefOr(attrs.DisplayName, "")), 
		Icon:        strings.TrimSpace(derefOr(attrs.Icon, "")),
		Args:        strings.TrimSpace(derefOr(attrs.Args, "")),
	}, diags
}



func parseCommand(block *hclsyntax.Block, c core.Common) (core.Command, hcl.Diagnostics) {


	var attrs struct {
		Exe  string  `hcl:"exe"`
		Args *string `hcl:"args,optional"`
	}


	if diags := gohcl.DecodeBody(block.Body, nil, &attrs); diags.HasErrors() {
		return core.Command{}, diags
	}

	var diags hcl.Diagnostics

	diags = append(diags, checkRequired(attrs.Exe, "exe", block.Body.Attributes["exe"].Expr.Range())...)
	diags = append(diags, checkOptional(attrs.Args, "args", attrRangeOf(block, "args"))...)



	if diags.HasErrors() {
		return core.Command{}, diags
	}

	return core.Command{
		Common: c,
		Exe:    strings.TrimSpace(attrs.Exe),
		Args:   strings.TrimSpace(derefOr(attrs.Args, "")),
	}, diags
}




func parseAddPath(block *hclsyntax.Block, c core.Common) (core.AddPath, hcl.Diagnostics) {
	var attrs struct {
		Dir string `hcl:"dir"`
	}
	if diags := gohcl.DecodeBody(block.Body, nil, &attrs); diags.HasErrors() {
		return core.AddPath{}, diags
	}

	var diags hcl.Diagnostics

	diags = append(diags, checkRequired(attrs.Dir, "dir", block.Body.Attributes["dir"].Expr.Range())...)

	if diags.HasErrors() {
		return core.AddPath{}, diags
	}

	return core.AddPath{Common: c, Dir: strings.TrimSpace(attrs.Dir)}, diags
}

// ---------- shared helpers ----------



func derefOr(p *string, def string) string {
	if p == nil {
		return def
	}
	return *p
}

func requiredErr(field string, rng hcl.Range) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  fmt.Sprintf("%s must not be empty", field),
		Subject:  rng.Ptr(),
	}
}




func checkRequired(value, field string, attrRange hcl.Range) hcl.Diagnostics {
	var diags hcl.Diagnostics

	if strings.TrimSpace(value) == "" {
		diags = append(diags, requiredErr(field, attrRange))
		return diags
	}
	return diags
}


func checkOptional(value *string, field string, attrRange hcl.Range) hcl.Diagnostics {
	var diags hcl.Diagnostics
	if value == nil {
		return diags
	}
	if strings.TrimSpace(*value) == "" {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("%s must not be empty if provided", field),
			Detail:   fmt.Sprintf("omit %s entirely to use the default, or provide a non-empty value", field),
			Subject:  attrRange.Ptr(),
		})
		return diags
	}

	return diags
}


func attrRangeOf(block *hclsyntax.Block, name string) hcl.Range {
	if attr, ok := block.Body.Attributes[name]; ok {
		return attr.Expr.Range()
	}
	return block.DefRange()
}











