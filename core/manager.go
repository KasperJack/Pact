package core



type Manager interface {
	Install(InstallArgs) error
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

