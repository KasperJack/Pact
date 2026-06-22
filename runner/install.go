package main

import (
	//"os"
	//"path"
	//"path/filepath"

	"github.com/kasperjack/pact/core"
	"github.com/kasperjack/pact/core/platform"
)

func install (pkg string, version string, arch platform.Arch) error {

	

	repo,err := NewLocalRepo("test-buckets")
	if err != nil {
		return err
	}

	localState := NewLocalState("installed") // local state should contain the lockFile ? 
	lockFile, err := NewLockFile("installed/lock.hcl")
	if err != nil {
		return err
	}

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