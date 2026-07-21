package manager

import (
		"github.com/kasperjack/pact/core"
		"os"
	//"path"
		"path/filepath"
)

type Manager interface {
	Install(args core.InstallArgs) error
}



func NewManager(localState core.LocalState, repo core.Repo,lf core.LockFile) Manager {

	exePath, _ := os.Executable()

	exeDir := filepath.Dir(exePath)


	return &pkgManager{
		localState: localState,
		repo:       repo,
		lockFile: lf,
		staging: NewStagingArea(filepath.Join(exeDir, "staging")),
	}

}

type pkgManager struct { 
	localState core.LocalState
	repo       core.Repo
	lockFile core.LockFile
	staging *StagingArea //EC: change to an interface 
}

