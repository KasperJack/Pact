package manager

import (
		"github.com/kasperjack/pact/core"
		"github.com/kasperjack/pact/core/repo"
		"github.com/kasperjack/pact/core/lockmanager"
)

type Manager interface {
	Install(args core.InstallArgs) error
}



func NewManager(localState *core.LocalState) (Manager,error) {


	localRepo, err := repo.NewLocalRepo(localState.Repo)
	if err != nil {
		return nil,err
	}



	userLockManger, err := lockmanager.New(localState.UserLockFile)
	if err != nil {
		return nil,err
	}

	systemLockManger, err := lockmanager.New(localState.SystemLockFile)
	if err != nil {
		return nil,err
	}

	

	return &pkgManager{
		localState: localState,
		repo:       localRepo,
		LockManagers: core.LockManagers{User: userLockManger,System: systemLockManger},
	},nil

}

type pkgManager struct { 
	localState *core.LocalState
	repo       core.Repo
	LockManagers core.LockManagers
}

