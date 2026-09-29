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
	"slices"
	// "stings" 💀
	//"/github.com/zclconf/go-cty/cty"
)

var (
	validIDPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

)

//TODO:
// Add table-driven tests
// Add a source composition layer (repo -> composition -> parser) 



var localBlockConstructors = map[string]func(string, hcl.Range) localBlock{
    core.Shortcut{}.Name(): func(id string, rng hcl.Range) localBlock {
        b := &shortcut{}
        b.setID(id)
        b.setRange(rng)
		b.setScope(core.Shortcut{}.SuppotedScopes())
        return b
    },
    core.AddPath{}.Name(): func(id string, rng hcl.Range) localBlock {
        b := &addPath{}
        b.setID(id)
        b.setRange(rng)
		b.setScope(core.AddPath{}.SuppotedScopes())
        return b
    },
    core.Command{}.Name(): func(id string, rng hcl.Range) localBlock {
        b := &command{}
        b.setID(id)
        b.setRange(rng)
		b.setScope(core.Command{}.SuppotedScopes())
        return b
    },
}


type localBlock interface {
	validate(*hclsyntax.Body) hcl.Diagnostics
	export() core.Block
	uniqueFields() map[string]string
	

	setID(string)
	setRange(hcl.Range)
	setScope([]core.Scope)

	

	getID() string
	getRange() hcl.Range
	supportedScopes() []core.Scope

}

type meta struct {
	ID    string
	Range hcl.Range
	Scope []core.Scope
}

func (m *meta) setID(id string)       { m.ID = id }
func (m *meta) getID() string         { return m.ID }
func (m *meta) setRange(r hcl.Range)  { m.Range = r }
func (m *meta) getRange() hcl.Range   { return m.Range }
func (m *meta) setScope(s []core.Scope) { m.Scope = s }
func (m *meta) supportedScopes() []core.Scope { return m.Scope }







// //////////// local types


type shortcut struct {
	meta

	DisplayName *string `hcl:"display_name,optional"`
	Exe         string  `hcl:"exe"`
	Icon        *string `hcl:"icon,optional"`
	Args        *string `hcl:"args,optional"`


}






func (s *shortcut) validate(body *hclsyntax.Body) hcl.Diagnostics {

	var diags hcl.Diagnostics


	rangeOf := func(field string) hcl.Range {

        return attrRangeOf(body, field)
    }


	//diags = append(diags, checkRequired(s.Exe, "exe", block.Body.Attributes["exe"].Expr.Range())...)

	diags = append(diags, checkRequired(s.Exe, "exe", rangeOf("exe"))...)
	diags = append(diags, checkOptional(s.DisplayName, "display_name", rangeOf("display_name"))...)
	diags = append(diags, checkOptional(s.Icon, "icon", rangeOf("icon"))...)
	diags = append(diags, checkOptional(s.Args, "args", rangeOf("args"))...)

	return diags

}


func (s *shortcut) export() core.Block {

		return core.Shortcut{
		ID:          s.ID,
		Exe:         strings.TrimSpace(s.Exe),
		DisplayName: strings.TrimSpace(derefOr(s.DisplayName, "")),
		Icon:        strings.TrimSpace(derefOr(s.Icon, "")),
		Args:        strings.TrimSpace(derefOr(s.Args, "")),
	}
}



func (s *shortcut) uniqueFields() map[string]string {
	m := map[string]string{}

	if v := strings.TrimSpace(derefOr(s.DisplayName, "")); v != "" {
		m["display_name"] = v
	}

	return m
}








type command struct {
	meta

	Exe  string  `hcl:"exe"`
	Args *string `hcl:"args,optional"`

}

func (c *command) validate(body *hclsyntax.Body) hcl.Diagnostics {

	var diags hcl.Diagnostics


	rangeOf := func(field string) hcl.Range {

        return attrRangeOf(body, field)
    }






	diags = append(diags, checkRequired(c.Exe, "exe", rangeOf("exe"))...)
	diags = append(diags, checkOptional(c.Args, "args", rangeOf("args"))...)

	return diags
}




func (c *command) export() core.Block {
		return core.Command{
		ID:   c.ID,
		Exe:  strings.TrimSpace(c.Exe),
		Args: strings.TrimSpace(derefOr(c.Args, "")),
	}
}

func (c *command) uniqueFields() map[string]string {
	return map[string]string{
		"exe": strings.TrimSpace(c.Exe),
	}
}









type addPath struct {
	meta
	Dir string `hcl:"dir"`

}



func (a *addPath) validate(body *hclsyntax.Body) hcl.Diagnostics {

	var diags hcl.Diagnostics

	rangeOf := func(field string) hcl.Range {

        return attrRangeOf(body, field)
    }


	diags = append(diags, checkRequired(a.Dir, "dir", rangeOf("dir"))...)


	return diags
}





func (a *addPath) export() core.Block {
	return core.AddPath{
		ID:  a.ID,
		Dir: strings.TrimSpace(a.Dir),
	}
}


func (a *addPath) uniqueFields() map[string]string {

	return nil
}















type registry struct {
	// blockType -> id -> block   (id uniqueness, only when id != "")
	ids map[string]map[string]localBlock

	// blockType -> fieldName -> value -> block   (field uniqueness within same type)
	fields map[string]map[string]map[string]localBlock
}

func newRegistry() *registry {
	return &registry{
		ids:    map[string]map[string]localBlock{},
		fields: map[string]map[string]map[string]localBlock{},
	}
}

func (r *registry) add(blockType string, b localBlock) hcl.Diagnostics {
	var diags hcl.Diagnostics
	rng := b.getRange()

	// --- ID uniqueness, only if id is present ---
	if id := b.getID(); id != "" {
		if r.ids[blockType] == nil {
			r.ids[blockType] = map[string]localBlock{}
		}
		if existing, exists := r.ids[blockType][id]; exists {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Duplicate ID",
				Detail: fmt.Sprintf("%s %q already defined at %s",
					blockType, id, existing.getRange()),
				Subject: rng.Ptr(),
			})
		} else {
			r.ids[blockType][id] = b
		}
	}

	// --- arbitrary field uniqueness, generic over whatever the type declares ---
	for field, val := range b.uniqueFields() {
		if r.fields[blockType] == nil {
			r.fields[blockType] = map[string]map[string]localBlock{}
		}
		if r.fields[blockType][field] == nil {
			r.fields[blockType][field] = map[string]localBlock{}
		}
		if existing, exists := r.fields[blockType][field][val]; exists {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("Duplicate %s", field),
				Detail: fmt.Sprintf("%s %q has the same %s (%q) as %s already defined at %s",
					blockType, b.getID(), field, val, blockType, existing.getRange()),
				Subject: rng.Ptr(),
			})
		} else {
			r.fields[blockType][field][val] = b
		}
	}

	return diags
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
	//check no top level attributes exist 


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



	resolved, d := parseScope(body,scope)
	diags = append(diags, d...)
	if diags.HasErrors() {
		return nil, diags
	}

	m := &core.Manifest{Scope: make(map[core.Scope]core.ManifestScope)}
	m.Scope[scope] = resolved

	return m, diags


}

// ---------- dual-scope manifest (user{}/system{} required) ----------

func parseDualScopeManifest(body *hclsyntax.Body, acceptScope []core.Scope) (*core.Manifest, hcl.Diagnostics) {

	m := &core.Manifest{Scope: map[core.Scope]core.ManifestScope{}}

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
			scope, d := parseScope(block.Body,core.ScopeUser)
			diags = append(diags, d...)
			m.Scope[core.ScopeUser] = scope

		case "system":
			if seenSystem {
				diags = append(diags, dupTopLevelErr("system", block.DefRange()))
				continue
			}
			seenSystem = true
			scope, d := parseScope(block.Body,core.ScopeSystem)
			diags = append(diags, d...)
			m.Scope[core.ScopeSystem] = scope



		default:

			if _, ok := localBlockConstructors[block.Type]; ok {
				diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("unexpected top-level %q block", block.Type),
				Detail:   "this package supports multiple scopes — wrap actions in user{} and/or system{}",
				Subject:  block.DefRange().Ptr(),
				})

				return nil, diags



			}else{

				diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("unknown top-level block type %q", block.Type),
				Detail:   `only "user" and "system" blocks are allowed at the top level`,
				Subject:  block.DefRange().Ptr(),
				})

				return nil, diags
			}


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

func parseScope(body *hclsyntax.Body, mScope core.Scope) (core.ManifestScope, hcl.Diagnostics) {

	var diags hcl.Diagnostics

	installPath, d := decodeInstallPath(body)
	diags = append(diags, d...)





	scope := &core.ManifestScope{InstallPath: installPath}

	reg := newRegistry()
	var parsed []localBlock

	for _, block := range body.Blocks {
		lb, d := parseBlock(block,mScope) // now returns a valid localBlock not core.Block
		diags = append(diags, d...)
		if lb == nil {
			continue
		}

		diags = append(diags, reg.add(block.Type, lb)...) // catches ID + field dupes
		parsed = append(parsed, lb)                        // single append point = order preserved
	}

	if diags.HasErrors() {
		return core.ManifestScope{}, diags
	}

	for _, lb := range parsed {
		scope.Blocks = append(scope.Blocks, lb.export()) // export only after everything's clean
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




func parseBlock(hclBlock *hclsyntax.Block, mScope core.Scope) (localBlock, hcl.Diagnostics) {
	id, diags := blockLabel(hclBlock)
	if diags.HasErrors() {
		return nil, diags
	}

	create, ok := localBlockConstructors[hclBlock.Type]
	if !ok {
		return nil, hcl.Diagnostics{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("unknown block type %q", hclBlock.Type),
			Detail:   `expected "shortcut", "command", or "add_path"`,
			Subject:  hclBlock.DefRange().Ptr(),
		}}
	}

	lb := create(id, hclBlock.DefRange())



	diags = gohcl.DecodeBody(hclBlock.Body, nil, lb)
	if diags.HasErrors() {
		return nil, diags
	}


	if !slices.Contains(lb.supportedScopes(), mScope) {
		return nil, hcl.Diagnostics{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("block %q does not support %s scope", hclBlock.Type, mScope),
			Detail:   fmt.Sprintf("%s scope is not supported for this block type", mScope),
			Subject:  hclBlock.DefRange().Ptr(),
		}}
	}


	diags = lb.validate(hclBlock.Body)
	if diags.HasErrors() {
		return nil, diags
	}

	return lb, diags
}





// blockLabel enforces the 0-or-1-label rule and validates a single label.
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













