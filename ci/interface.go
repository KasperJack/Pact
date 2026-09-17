package main


import (
	"github.com/kasperjack/pact/core"
	"github.com/hashicorp/hcl/v2"
)

type ValidatedInterface struct {
	User   []core.Option
	System []core.Option
}




func ValidateInterface(m *core.Interface) (*ValidatedInterface, hcl.Diagnostics) {}