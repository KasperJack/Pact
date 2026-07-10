package manager

import (
		"github.com/kasperjack/pact/core"
)

type Manager interface {
	Install(InstallArgs) error
}

func NewManager(localState core.LocalState, repo core.Repo,lf core.LockFile) Manager {

	return &pkgManager{
		localState: localState,
		repo:       repo,
		lockFile: lf,
		staging: NewStagingArea("C:\\Users\\Aya\\Desktop\\pact\\bin\\staging"),
	}

}

type pkgManager struct { 
	localState core.LocalState
	repo       core.Repo
	lockFile core.LockFile
	staging *StagingArea //EC: change to an interface 
}

