package main

import (
	"os"
	//"path"
	"path/filepath"

	"github.com/kasperjack/pact/core"
	"github.com/kasperjack/pact/core/manager"
	"github.com/kasperjack/pact/core/repo"
	"github.com/kasperjack/pact/core/lockmanager"
)





func install(pkg string, version string, arch core.Arch) error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exeDir := filepath.Dir(exePath)

	localState := NewLocalState(
		filepath.Join(exeDir, "desktop"),
		filepath.Join(exeDir, "cache"),
		filepath.Join(exeDir, "pkg"),
		filepath.Join(exeDir, "repo"),
		filepath.Join(exeDir, "lock.hcl"),
	)

	localRepo, err := repo.NewLocalRepo(localState.Repo())
	if err != nil {
		return err
	}

	lm, err := lockmanager.New(localState.LockFile())
	if err != nil {
		return err
	}

	m := manager.NewManager(localState, localRepo, lm)

	err = m.Install(core.InstallArgs{
		PackageIdentifier: pkg,
		Version:           core.ParseVersion(version),
		TargetArch:        arch,
	})

	
	if err != nil {
		return err
	}

	return nil
}
