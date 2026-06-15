package main

import (
	//"os"
	//"path"
	//"path/filepath"

	"github.com/kasperjack/pact/core"
	"github.com/kasperjack/pact/core/platform"
)

func install (pkg string, version string, arch platform.Arch) error {

	

	repo := NewLocalRepo("test-buckets")
	localState := NewLocalState("C:\\Users\\Aya\\Desktop\\pact\\bin\\installed")
	lockFile, err := NewLockFile("installed/lock.hcl")


	m := core.NewManager(localState,repo,lockFile)


	err = m.Install(core.InstallArgs{Name: pkg,Version: version,Arch: arch}) // move InstallArgs to core 
	if err != nil {
		return err
	}

	

	return nil
}

//TODO: 
// lib will take 3 manifest of the pcakge + install location that the pkmg uses 
// it will retrun an interface install update remove 