package main

import (
	//"os"
	//"path"
	//"path/filepath"

	//"github.com/kasperjack/pact/core"
	"github.com/kasperjack/pact/core/manager"
)

func install (pkg string, version string, arch string) error {

	

	repo,err := NewLocalRepo("test-buckets")
	if err != nil {
		return err
	}

	localState := NewLocalState("installed") // local state should contain the lockFile ? 
	lockFile, err := NewLockFile("installed/lock.hcl")
	if err != nil {
		return err
	}

	m := manager.NewManager(localState,repo,lockFile)


	err = m.Install(manager.InstallArgs{Name: pkg,Version: version,TargetArch: arch}) // move InstallArgs to core 
	if err != nil {
		return err
	}

	

	return nil
}

//TODO: 
// lib will take 3 manifest of the pcakge + install location that the pkmg uses 
// it will retrun an interface install update remove 