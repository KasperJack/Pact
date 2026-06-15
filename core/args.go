package core

import (
	
	"github.com/kasperjack/pact/core/platform"
)
type InstallArgs struct {
    Name       string
    Version    string
	Arch  platform.Arch
}