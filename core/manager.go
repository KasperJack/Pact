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

	// check if pcakge is already installed 
	// ehck if pakge supports multiple version 


	ok, err := m.repo.PackageExists(pkgName)

	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("package not found")
	}

	


	pf,err := m.repo.LoadPackage(pkgName,"2.6.1")
	if err != nil {
		return err
	}


	//fmt.Println(pf.Package.Description)
	//fmt.Println(pf.Release.Source.X86.URL)



	return nil
}
