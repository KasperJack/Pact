package manager

import (
		"github.com/kasperjack/pact/core"
)

type Manager interface {
	Install(args core.InstallArgs) error
}



func NewManager(localState core.LocalState, repo core.Repo,lm core.LockManager) Manager {


	return &pkgManager{
		localState: localState,
		repo:       repo,
		lockM: lm,
	}

}

type pkgManager struct { 
	localState core.LocalState
	repo       core.Repo
	lockM core.LockManager
}

