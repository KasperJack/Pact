package parce

import (
	"fmt"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/kasperjack/pact/core"
)






func Manifest(path string) (*core.Manifest, error) {
	parser := hclparse.NewParser()
	f, diags := parser.ParseHCLFile(path)
	if diags.HasErrors() {
		return nil, fmt.Errorf("parsing %s: %w", path, diags)
	}

	syntaxBody, ok := f.Body.(*hclsyntax.Body)
	if !ok {
		return nil, fmt.Errorf("unexpected body type for %s", path)
	}

	m := core.NewManifest()

	for _, block := range syntaxBody.Blocks {
		switch block.Type {
		case "shortcut":
			var sb core.ShortcutBody
			if diags := gohcl.DecodeBody(block.Body, nil, &sb); diags.HasErrors() {
				return nil, fmt.Errorf("decoding shortcut block: %w", diags)
			}
			if err := addByLabel(block.Labels, sb, m.Shortcuts); err != nil {
				return nil, fmt.Errorf("shortcut: %w", err)
			}

		case "command":
			var cb core.CommandBody
			if diags := gohcl.DecodeBody(block.Body, nil, &cb); diags.HasErrors() {
				return nil, fmt.Errorf("decoding command block: %w", diags)
			}
			if err := addByLabel(block.Labels, cb, m.Commands); err != nil {
				return nil, fmt.Errorf("command: %w", err)
			}



		case "scope":
			if m.Scope != nil {
				return nil, fmt.Errorf("manifest declares more than one scope block")
			}
			var sc core.Scope
			if diags := gohcl.DecodeBody(block.Body, nil, &sc); diags.HasErrors() {
				return nil, fmt.Errorf("decoding scope block: %w", diags)
			}
			m.Scope = &sc



		default:
			return nil, fmt.Errorf("unknown block type %q", block.Type)
		}
	}

	if err := validateScope(m.Scope); err != nil {
		return nil, err
	}

	return m, nil
}





// addByLabel handles the shared 0-or-1-label rule for any action type.
func addByLabel[T any](labels []string, body T, set *core.ActionSet[T]) error {
	switch len(labels) {
	case 0:
		set.AddUnconditional(body)
		return nil
	case 1:
		return set.AddTagged(labels[0], body)
	default:
		return fmt.Errorf("block takes 0 or 1 labels, got %d", len(labels))
	}
}



func validateScope(s *core.Scope) error {
	if s == nil {
		return fmt.Errorf("manifest must declare a scope block")
	}

	if s.User == nil && s.System == nil {
		return fmt.Errorf("scope must declare exactly one of user or system")
	}
	return nil
}



