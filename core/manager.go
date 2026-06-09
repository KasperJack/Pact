package core

import "github.com/kasperjack/pact/core/model"

type Manager interface {
	Install(model.InstallArgs) error
}

func NewManager(localState LocalState, repo Repo,lf LockFile) Manager {

	return &pkgManager{
		localState: localState,
		repo:       repo,
		lockFile: lf,
	}

}

type pkgManager struct { 
	localState LocalState
	repo       Repo
	lockFile LockFile
}

