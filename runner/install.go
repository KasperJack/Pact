package main

import (
	"os"
	//"path"
	"path/filepath"

	"github.com/kasperjack/pact/core"
	"github.com/kasperjack/pact/core/manager"
)

func install (pkg string, version string, arch string) error {

	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exeDir := filepath.Dir(exePath)


	repo,err := NewLocalRepo(filepath.Join(exeDir, "test-buckets"))
	if err != nil {
		return err
	}

	localState := NewLocalState(filepath.Join(exeDir, "installed")) // local state should contain the lockFile ? 
	lockFile, err := NewLockFile(filepath.Join(exeDir, "installed", "lock.hcl"))
	if err != nil {
		return err
	}

	m := manager.NewManager(localState,repo,lockFile)


	err = m.Install(core.InstallArgs{PackageIdentifier: pkg,Version: version,TargetArch: arch}) // move InstallArgs to core 
	if err != nil {
		return err
	}

	

	return nil
}

//TODO: 
// lib will take 3 manifest of the pcakge + install location that the pkmg uses 
// it will retrun an interface install update remove 