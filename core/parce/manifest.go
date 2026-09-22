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
	//"/github.com/zclconf/go-cty/cty"
)

var (
	validIDPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

	//numbers = [5]core.Block{}
)
/*
var m = map[string]localBlock{
    core.Shortcut{}.Name(): &shortcut{},
    core.AddPath{}.Name():  &addPath{},
    core.Command{}.Name():  &command{},
} // using reflect ? 
*/

var localBlockConstructors = map[string]func(string, hcl.Range) localBlock{
    core.Shortcut{}.Name(): func(id string, rng hcl.Range) localBlock {
        b := &shortcut{}
        b.setID(id)
        b.setRange(rng)
        return b
    },
    core.AddPath{}.Name(): func(id string, rng hcl.Range) localBlock {
        b := &addPath{}
        b.setID(id)
        b.setRange(rng)
        return b
    },
    core.Command{}.Name(): func(id string, rng hcl.Range) localBlock {
        b := &command{}
        b.setID(id)
        b.setRange(rng)
        return b
    },
}


type localBlock interface {
	validate(*hclsyntax.Block) hcl.Diagnostics
	export() core.Block

	

	setID(string)
	getID() string
	setRange(hcl.Range)
	getRange() hcl.Range

}




// //////////// local types


type shortcut struct {
	ID string

	DisplayName *string `hcl:"display_name,optional"`
	Exe         string  `hcl:"exe"`
	Icon        *string `hcl:"icon,optional"`
	Args        *string `hcl:"args,optional"`

	Range hcl.Range
}






func (s *shortcut) validate(block *hclsyntax.Block) hcl.Diagnostics {
	return nil
}
func (s *shortcut) export() core.Block {
	return nil
}
func (s *shortcut) setID(id string) {
	s.ID = id
}
func (s *shortcut) setRange(r hcl.Range) {
	s.Range = r
}








type command struct {
	ID   string
	Exe  string  `hcl:"exe"`
	Args *string `hcl:"args,optional"`

	Range hcl.Range
}

func (c *command) validate(block *hclsyntax.Block) hcl.Diagnostics {
	return nil
}
func (c *command) export() core.Block {
	return nil
}


func (c *command) uniqueFields() map[string]string {

	m := map[string]string{}

	m["exe"] = strings.TrimSpace(c.Exe)


	
	return m
}






func (c *command) setID(id string) {
	c.ID = id
}
func (c *command) setRange(r hcl.Range) {
	c.Range = r
}













type addPath struct {
	ID string

	Dir string `hcl:"dir"`

	Range hcl.Range
}

func (a *addPath) validate(block *hclsyntax.Block) hcl.Diagnostics {

	var diags hcl.Diagnostics

	diags = append(diags, checkRequired(a.Dir, "dir", block.Body.Attributes["dir"].Expr.Range())...)


	return diags
}



func (a *addPath) export() core.Block {
	return core.AddPath{
		ID:  a.ID,
		Dir: strings.TrimSpace(a.Dir),
	}
}


func (a *addPath) setID(id string) {
	a.ID = id
}
func (a *addPath) setRange(r hcl.Range) {
	a.Range = r
}











func Manifest(src []byte, acceptScope []core.Scope) (*core.Manifest, hcl.Diagnostics) {
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
			Detail:   "could not parse manifest.hcl as HCL native syntax",
		}}
	}

	if len(acceptScope) == 1 {
		return parseSingleScopeManifest(syntaxBody, acceptScope[0])
	}

	return parseDualScopeManifest(syntaxBody, acceptScope)
}

// ---------- single-scope manifest (no wrapper) ----------

func parseSingleScopeManifest(body *hclsyntax.Body, scope core.Scope) (*core.Manifest, hcl.Diagnostics) {

	var diags hcl.Diagnostics

	for _, block := range body.Blocks { // top level blocks check

		if _, ok := localBlockConstructors[block.Type]; ok {
			// handled below via parseScope
			continue
		}	


		switch block.Type {

		case "user", "system":
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("unexpected %q block", block.Type),
				Detail:   "this package only supports one scope — declare install_path and actions directly at the top level, without a user{}/system{} wrapper",
				Subject:  block.DefRange().Ptr(),
			})
		default:
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("unknown top-level block type %q", block.Type),
				Detail:   `only "shortcut", "command", or "add_path" are allowed at the top level for a single-scope package`,
				Subject:  block.DefRange().Ptr(),
			})
		}
	}


	if diags.HasErrors() {
		return nil, diags
	}



	resolved, d := parseScope(body)
	diags = append(diags, d...)
	if diags.HasErrors() {
		return nil, diags
	}

	m := &core.Manifest{
    Scope: make(map[core.Scope]core.ManifestScope),
	}

	switch scope {
	case core.ScopeUser:
		m.Scope[core.ScopeUser] = resolved
	case core.ScopeSystem:
		m.Scope[core.ScopeSystem] = resolved
	default:
		return nil, hcl.Diagnostics{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("unknown scope %v", scope),
		}}
	}

	return m, diags
}

// ---------- dual-scope manifest (user{}/system{} required) ----------

func parseDualScopeManifest(body *hclsyntax.Body, acceptScope []core.Scope) (*core.Manifest, hcl.Diagnostics) {
	m := &core.Manifest{}
	var diags hcl.Diagnostics
	seenUser := false
	seenSystem := false

	for _, block := range body.Blocks {
		switch block.Type {

		case "user":
			if seenUser {
				diags = append(diags, dupTopLevelErr("user", block.DefRange()))
				continue
			}
			seenUser = true
			scope, d := parseScope(block.Body)
			diags = append(diags, d...)
			m.Scope[core.ScopeUser] = scope

		case "system":
			if seenSystem {
				diags = append(diags, dupTopLevelErr("system", block.DefRange()))
				continue
			}
			seenSystem = true
			scope, d := parseScope(block.Body)
			diags = append(diags, d...)
			m.Scope[core.ScopeSystem] = scope

		case "shortcut", "command", "add_path":
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("unexpected top-level %q block", block.Type),
				Detail:   "this package supports multiple scopes — wrap actions in user{} and/or system{}",
				Subject:  block.DefRange().Ptr(),
			})

		default:
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("unknown top-level block type %q", block.Type),
				Detail:   `only "user" and "system" blocks are allowed at the top level`,
				Subject:  block.DefRange().Ptr(),
			})
		}
	}

	// require a wrapper for every scope this package is declared to accept —
	// missing one silently would leave that scope with no install_path/actions
	for _, s := range acceptScope {
		switch s {
		case core.ScopeUser:
			if !seenUser {
				diags = append(diags, missingScopeBlockErr("user", body))
			}
		case core.ScopeSystem:
			if !seenSystem {
				diags = append(diags, missingScopeBlockErr("system", body))
			}
		}
	}

	if diags.HasErrors() {
		return nil, diags
	}
	return m, diags
}

func missingScopeBlockErr(name string, body *hclsyntax.Body) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  fmt.Sprintf("missing %q block", name),
		Detail:   fmt.Sprintf("package accepts %q scope, but manifest.hcl has no %s{} block", name, name),
		Subject:  body.SrcRange.Ptr(),
	}
}

func dupTopLevelErr(blockType string, rng hcl.Range) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  fmt.Sprintf("duplicate %s block", blockType),
		Detail:   fmt.Sprintf("a manifest may declare at most one %s block", blockType),
		Subject:  rng.Ptr(),
	}
}

// ---------- scope-level parse (shared by both paths) ----------

func parseScope(body *hclsyntax.Body) (core.ManifestScope, hcl.Diagnostics) {

	var diags hcl.Diagnostics

	installPath, d := decodeInstallPath(body)
	diags = append(diags, d...)

	scope := &core.ManifestScope{InstallPath: installPath}

	for _, block := range body.Blocks {
		b, d := parseBlock(block)
		diags = append(diags, d...)
		if b != nil {
			scope.Blocks = append(scope.Blocks, b) // single append point = order preserved
		}
	}

	return *scope, diags
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
				// change name to hclBlock
func parseBlock(block *hclsyntax.Block) (core.Block, hcl.Diagnostics) {

	id, diags := blockLabel(block)

	if diags.HasErrors() {
		return nil, diags
	}

	create, ok := localBlockConstructors[block.Type]

	if !ok {
		return nil, hcl.Diagnostics{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("unknown block type %q", block.Type),
			Detail:   `expected "shortcut", "command", or "add_path"`,
			Subject:  block.DefRange().Ptr(),
		}}
	}

	localBlock := create(id,block.DefRange())


	diags = gohcl.DecodeBody(block.Body, nil, localBlock)
	if diags.HasErrors() {
		return nil, diags
	}

	//localBlock.setID(id)
	//localBlock.setRange(block.DefRange())


	diags = localBlock.validate(block)
	if diags.HasErrors() {
		return nil, diags
	}




	return localBlock.export(), diags

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

func parseShortcut(block *hclsyntax.Block, id string) (core.Shortcut, hcl.Diagnostics) {
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
		ID:          id,
		Exe:         strings.TrimSpace(attrs.Exe),
		DisplayName: strings.TrimSpace(derefOr(attrs.DisplayName, "")),
		Icon:        strings.TrimSpace(derefOr(attrs.Icon, "")),
		Args:        strings.TrimSpace(derefOr(attrs.Args, "")),
	}, diags
}

func parseCommand(block *hclsyntax.Block, id string) (core.Command, hcl.Diagnostics) {

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
		ID:   id,
		Exe:  strings.TrimSpace(attrs.Exe),
		Args: strings.TrimSpace(derefOr(attrs.Args, "")),
	}, diags
}

func parseAddPath(block *hclsyntax.Block, id string) (core.AddPath, hcl.Diagnostics) {


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

	return core.AddPath{ID: id, Dir: strings.TrimSpace(attrs.Dir)}, diags
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
