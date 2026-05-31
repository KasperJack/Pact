package main

import (
	//"os"
	//"path"
	//"path/filepath"

	"github.com/kasperjack/pact/core"
)

func install (pkg string, version string) error {

	//r := repo{}

	repo := NewLocalRepo("C:\\Users\\Aya\\Desktop\\pact\\bin\\test-buckets")
	localState := NewLocalState("ass")
	

	m := core.NewManager(localState,repo)


	err := m.Install(pkg)
	if err != nil {
		return err
	}

	return nil
}


// lib will take 3 manifest of the pcakge + install location that the pkmg uses 
// it will retrun an interface install update remove 