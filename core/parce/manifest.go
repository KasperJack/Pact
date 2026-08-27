package parce

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/gohcl"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/kasperjack/pact/core"
	"github.com/zclconf/go-cty/cty"
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

		case "user":
			if m.User != nil {
				return nil, fmt.Errorf("manifest declares more than one user block")
			}
			mb, err := parseScopedBody(block.Body)
			if err != nil {
				return nil, fmt.Errorf("user: %w", err)
			}
			m.User = mb

		case "system":
			if m.System != nil {
				return nil, fmt.Errorf("manifest declares more than one system block")
			}
			mb, err := parseScopedBody(block.Body)
			if err != nil {
				return nil, fmt.Errorf("system: %w", err)
			}
			m.System = mb

		default:
			return nil, fmt.Errorf("unknown block type %q", block.Type)
		}
	}

	if m.User == nil && m.System == nil {
		return nil, fmt.Errorf("manifest must declare at least one of user or system")
	}

	return m, nil
}

// parseScopedBody handles a single user{} or system{} block: its
// install_path attribute plus its nested shortcut/command/add_path blocks.
func parseScopedBody(body *hclsyntax.Body) (*core.ManifestBody, error) {
	installPath, err := decodeInstallPath(body)
	if err != nil {
		return nil, err
	}

	mb := core.NewManifestBody()
	mb.InstallPath = installPath

	for _, block := range body.Blocks {
		if err := parseOneBlock(block, mb); err != nil {
			return nil, err
		}
	}

	return mb, nil
}

// decodeInstallPath pulls the required install_path attribute directly off
// the block body. We do this by hand rather than via gohcl.DecodeBody on the
// whole body, since that body also contains the nested shortcut/command/
// add_path blocks, which we're walking manually below.
func decodeInstallPath(body *hclsyntax.Body) (string, error) {
	attr, ok := body.Attributes["install_path"]
	if !ok {
		return "", fmt.Errorf("missing required attribute install_path")
	}

	val, diags := attr.Expr.Value(nil)
	if diags.HasErrors() {
		return "", fmt.Errorf("evaluating install_path: %w", diags)
	}
	if val.Type() != cty.String {
		return "", fmt.Errorf("install_path must be a string")
	}

	return val.AsString(), nil
}

// parseOneBlock decodes a single shortcut/command/add_path block into mb.
func parseOneBlock(block *hclsyntax.Block, mb *core.ManifestBody) error {
	switch block.Type {
	case "shortcut":
		var sb core.ShortcutBody
		if diags := gohcl.DecodeBody(block.Body, nil, &sb); diags.HasErrors() {
			return fmt.Errorf("decoding shortcut block: %w", diags)
		}

		if err := addByLabel(block.Labels, sb, mb.Shortcuts); err != nil {
			return fmt.Errorf("shortcut: %w", err)
		}

	case "command":
		var cb core.CommandBody
		if diags := gohcl.DecodeBody(block.Body, nil, &cb); diags.HasErrors() {
			return fmt.Errorf("decoding command block: %w", diags)
		}
		if err := addByLabel(block.Labels, cb, mb.Commands); err != nil {
			return fmt.Errorf("command: %w", err)
		}

	case "add_path":
		var ab core.AddPathBody
		if diags := gohcl.DecodeBody(block.Body, nil, &ab); diags.HasErrors() {
			return fmt.Errorf("decoding add_path block: %w", diags)
		}
		if err := addByLabel(block.Labels, ab, mb.AddPaths); err != nil {
			return fmt.Errorf("add_path: %w", err)
		}

	default:
		return fmt.Errorf("unknown block type %q inside user/system block", block.Type)
	}
	return nil
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


