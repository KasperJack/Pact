package main

import (
	//"os"
	//"fmt"
	//"reflect"
	
	"github.com/kasperjack/pact/core"
	"github.com/kasperjack/pact/core/manager"

)





func install(pkg string, version string, arch core.Arch) error {

	localState := NewLocalState()



	m,err := manager.NewManager(localState)
	if err != nil {
		return err
	}



	
	err = m.Install(core.InstallArgs{
		PackageIdentifier: pkg,
		Version:           core.ParseVersion(version),
		TargetArch:        arch,
		Scope:             core.ScopeUndefined,
	})

	
	if err != nil {
		return err
	}

	return nil
}


