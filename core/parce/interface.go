package parce

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/kasperjack/pact/core"
	"github.com/hashicorp/hcl/v2/gohcl"

)

func Interface(path string) (*core.Interface, error) {
	parser := hclparse.NewParser()
	f, diags := parser.ParseHCLFile(path)
	if diags.HasErrors() {
		return nil, fmt.Errorf("parsing %s: %w", path, diags)
	}

	syntaxBody, ok := f.Body.(*hclsyntax.Body)
	if !ok {
		return nil, fmt.Errorf("unexpected body type for %s", path)
	}

	iface := core.NewInterface()

	for _, block := range syntaxBody.Blocks {
		switch block.Type {

		case "user":
			if iface.User != nil {
				return nil, fmt.Errorf("interface declares more than one user block")
			}
			ib, err := parseInterfaceBody(block.Body)
			if err != nil {
				return nil, fmt.Errorf("user: %w", err)
			}
			iface.User = ib

		case "system":
			if iface.System != nil {
				return nil, fmt.Errorf("interface declares more than one system block")
			}
			ib, err := parseInterfaceBody(block.Body)
			if err != nil {
				return nil, fmt.Errorf("system: %w", err)
			}
			iface.System = ib

		default:
			return nil, fmt.Errorf("unknown block type %q", block.Type)
		}
	}

	if iface.User == nil && iface.System == nil {
		return nil, fmt.Errorf("interface must declare at least one of user or system")
	}

	return iface, nil
}

func parseInterfaceBody(body *hclsyntax.Body) (*core.InterfaceBody, error) {
	ib := core.NewInterfaceBody()

	for _, block := range body.Blocks {
		if block.Type != "option" {
			return nil, fmt.Errorf("unknown block type %q inside user/system block", block.Type)
		}
		if len(block.Labels) != 1 {
			return nil, fmt.Errorf("option block takes exactly 1 label, got %d", len(block.Labels))
		}
		id := block.Labels[0]
		if _, exists := ib.Options[id]; exists {
			return nil, fmt.Errorf("duplicate option %q", id)
		}

		opt, err := parseOption(block.Body)
		if err != nil {
			return nil, fmt.Errorf("option %q: %w", id, err)
		}
		ib.Options[id] = opt
	}

	return ib, nil
}

func parseOption(body *hclsyntax.Body) (core.OptionBody, error) {
	var opt core.OptionBody
	if diags := gohcl.DecodeBody(body, nil, &opt); diags.HasErrors() {
		return opt, fmt.Errorf("decoding option block: %w", diags)
	}
	if len(opt.Binding) == 0 {
		return opt, fmt.Errorf("binding must contain at least one entry")
	}
	return opt, nil
}