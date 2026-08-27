package main

import (
	"fmt"
	"strings"
	"log"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"regexp"
)



var (
	validIDPattern    = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	invalidPathChars  = regexp.MustCompile(`[<>"|?*]`)
)


// ---------- domain types (no HCL knowledge, no pointers) ----------

type ResolvedBlock interface {
	blockID() string
	blockRange() hcl.Range
	isResolvedBlock()
}

type common struct {
	ID    string
	Range hcl.Range
}

func (c common) blockID() string       { return c.ID }
func (c common) blockRange() hcl.Range { return c.Range }
func (c common) isResolvedBlock()      {}

type Shortcut struct {
	common
	DisplayName string //optianl // // trim tarling and ending spacse 
	Exe         string // required // trim tarling and ending spacse  //check in path 
	Icon        string //optianl // trim tarling and ending spacse    //check in path
	Args        string //optianl // // trim tarling and ending spacse 
}

type Command struct {
	common
	Exe  string  // required // trim tarling and ending spacse  //check in path
	Args string //optianl // // trim tarling and ending spacse 
}

type AddPath struct {
	common
	Dir string // required // trim tarling and ending spacse  //check in path
}

type ResolvedScope struct {
	InstallPath string
	Blocks      []ResolvedBlock // all types, file order preserved
}

type ResolvedManifest struct {
	User   *ResolvedScope
	System *ResolvedScope
}






// ---------- top-level parse ----------

func ParseManifest(path string) (*ResolvedManifest, hcl.Diagnostics) {
	parser := hclparse.NewParser()
	f, diags := parser.ParseHCLFile(path)
	if diags.HasErrors() {
		return nil, diags
	}

	syntaxBody, ok := f.Body.(*hclsyntax.Body)
	if !ok {
		return nil, hcl.Diagnostics{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "unexpected body type",
			Detail:   fmt.Sprintf("could not parse %s as HCL native syntax", path),
		}}
	}

	m := &ResolvedManifest{}

	for _, block := range syntaxBody.Blocks {
		switch block.Type {


		case "user":
			if m.User != nil {
				diags = append(diags, dupTopLevelErr("user", block.DefRange()))
				continue
			}
			scope, d := parseScope(block.Body)
			diags = append(diags, d...)
			m.User = scope



		case "system":
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

func parseScope(body *hclsyntax.Body) (*ResolvedScope, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	installPath, d := decodeInstallPath(body)
	diags = append(diags, d...)

	scope := &ResolvedScope{InstallPath: installPath}

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

func parseBlock(block *hclsyntax.Block) (ResolvedBlock, hcl.Diagnostics) {

	id, diags := blockLabel(block)

	if diags.HasErrors() {
		return nil, diags
	}

	c := common{ID: id, Range: block.DefRange()}

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




func parseShortcut(block *hclsyntax.Block, c common) (Shortcut, hcl.Diagnostics) {
	var attrs struct {
		DisplayName *string `hcl:"display_name,optional"`
		Exe         string  `hcl:"exe"`
		Icon        *string `hcl:"icon,optional"`
		Args        *string `hcl:"args,optional"`
	}


	if diags := gohcl.DecodeBody(block.Body, nil, &attrs); diags.HasErrors() {
		return Shortcut{}, diags
	}


	var diags hcl.Diagnostics

	diags = append(diags, checkRequired(attrs.Exe, "exe", block.Body.Attributes["exe"].Expr.Range())...)
	diags = append(diags, checkOptional(attrs.DisplayName, "display_name", attrRangeOf(block, "display_name"))...)
	diags = append(diags, checkOptional(attrs.Icon, "icon", attrRangeOf(block, "icon"))...)
	diags = append(diags, checkOptional(attrs.Args, "args", attrRangeOf(block, "args"))...)


	if diags.HasErrors() {
		return Shortcut{}, diags
	}

	return Shortcut{
		common:      c,
		Exe:         strings.TrimSpace(attrs.Exe),
		DisplayName: strings.TrimSpace(derefOr(attrs.DisplayName, "")), 
		Icon:        strings.TrimSpace(derefOr(attrs.Icon, "")),
		Args:        strings.TrimSpace(derefOr(attrs.Args, "")),
	}, diags
}



func parseCommand(block *hclsyntax.Block, c common) (Command, hcl.Diagnostics) {


	var attrs struct {
		Exe  string  `hcl:"exe"`
		Args *string `hcl:"args,optional"`
	}


	if diags := gohcl.DecodeBody(block.Body, nil, &attrs); diags.HasErrors() {
		return Command{}, diags
	}

	var diags hcl.Diagnostics

	diags = append(diags, checkRequired(attrs.Exe, "exe", block.Body.Attributes["exe"].Expr.Range())...)
	diags = append(diags, checkOptional(attrs.Args, "args", attrRangeOf(block, "args"))...)



	if diags.HasErrors() {
		return Command{}, diags
	}

	return Command{
		common: c,
		Exe:    strings.TrimSpace(attrs.Exe),
		Args:   strings.TrimSpace(derefOr(attrs.Args, "")),
	}, diags
}




func parseAddPath(block *hclsyntax.Block, c common) (AddPath, hcl.Diagnostics) {
	var attrs struct {
		Dir string `hcl:"dir"`
	}
	if diags := gohcl.DecodeBody(block.Body, nil, &attrs); diags.HasErrors() {
		return AddPath{}, diags
	}

	var diags hcl.Diagnostics

	diags = append(diags, checkRequired(attrs.Dir, "dir", block.Body.Attributes["dir"].Expr.Range())...)

	if diags.HasErrors() {
		return AddPath{}, diags
	}

	return AddPath{common: c, Dir: strings.TrimSpace(attrs.Dir)}, diags
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

































type ValidatedManifest struct {
	User   *ResolvedScope
	System *ResolvedScope
}




func Validate(m *ResolvedManifest) (*ValidatedManifest, hcl.Diagnostics) {


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




func validateScope(scope *ResolvedScope) hcl.Diagnostics {
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
		case Shortcut:
			diags = append(diags, validateShortcut(v, shortcutIDs, shortcutExes)...)
		case Command:
			diags = append(diags, validateCommand(v, commandIDs,commandExes)...)
		case AddPath:
			diags = append(diags, validateAddPath(v, addPathIDs,addPathDirs)...)
		}
	}

	return diags
}





func validateShortcut(s Shortcut, ids, exes map[string]hcl.Range) hcl.Diagnostics {
	var diags hcl.Diagnostics

	diags = append(diags, checkDup(ids, s.ID, s.Range, "shortcut id")...)
	diags = append(diags, checkDup(exes, s.Exe, s.Range, "shortcut exe")...)
	diags = append(diags, checkValidPath(s.Exe, "exe", s.Range)...)
	diags = append(diags, checkValidPath(s.Icon, "icon", s.Range)...)

	return diags
}

func validateCommand(c Command, ids, exes map[string]hcl.Range) hcl.Diagnostics {
	var diags hcl.Diagnostics

	diags = append(diags, checkDup(ids, c.ID, c.Range, "command id")...)
	diags = append(diags, checkDup(exes, c.Exe, c.Range, "command exe")...)
	diags = append(diags, checkValidPath(c.Exe, "exe", c.Range)...)

	return diags
}

func validateAddPath(a AddPath, ids, dirs map[string]hcl.Range) hcl.Diagnostics {
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














func main() {
	m, diags := ParseManifest("mfs.hcl")
	if diags.HasErrors() {
		printDiags(diags)
		log.Fatal("parse failed")
	}

	vm, diags := Validate(m)
	if diags.HasErrors() {
		printDiags(diags)
		log.Fatal("validation failed")
	}

	printScope("user", vm.User)
	printScope("system", vm.System)
}

func printScope(name string, scope *ResolvedScope) {
	if scope == nil {
		fmt.Printf("%s: (not declared)\n", name)
		return
	}

	fmt.Printf("%s: install_path=%q\n", name, scope.InstallPath)
	for i, b := range scope.Blocks {
		switch v := b.(type) {
		case Shortcut:
			fmt.Printf("  [%d] shortcut id=%q display_name=%q exe=%q icon=%q args=%q\n",
				i, v.ID, v.DisplayName, v.Exe, v.Icon, v.Args)
		case Command:
			fmt.Printf("  [%d] command id=%q exe=%q args=%q\n",
				i, v.ID, v.Exe, v.Args)
		case AddPath:
			fmt.Printf("  [%d] add_path id=%q dir=%q\n",
				i, v.ID, v.Dir)
		}
	}
}

func printDiags(diags hcl.Diagnostics) {
	for _, d := range diags {
		fmt.Println(d.Error())
	}
}