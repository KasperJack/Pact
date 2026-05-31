package core

import "fmt"

type Manager interface {
	Install(string) error
}

func NewManager(localState LocalState, repo Repo) Manager {

	return &pkgManager{
		localState: localState,
		repo:       repo,
	}

}

type pkgManager struct { // lowercase, hidden from outside
	localState LocalState
	repo       Repo
}

func (m *pkgManager) Install(pkgName string) error {

	ok, err := m.repo.PackageExists(pkgName)

	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("package not found")
	}

	


	return nil
}
