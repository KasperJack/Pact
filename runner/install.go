package main

import (
	//"os"
	//"path"
	//"path/filepath"

	"github.com/kasperjack/pact/core"
	"github.com/kasperjack/pact/core/model"
)

func install (pkg string, version string) error {

	

	repo := NewLocalRepo("test-buckets")
	localState := NewLocalState("ass")
	lockFile, err := NewLockFile("installed/lock.hcl")


	m := core.NewManager(localState,repo,lockFile)


	err = m.Install(model.InstallArgs{})
	if err != nil {
		return err
	}



	return nil
}


// lib will take 3 manifest of the pcakge + install location that the pkmg uses 
// it will retrun an interface install update remove 