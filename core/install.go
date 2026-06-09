package core
import (
	"fmt"
	"github.com/kasperjack/pact/core/model"
)




func (m *pkgManager) Install(agrs model.InstallArgs) error {

	_ , err := m.lockFile.GetInstalled(agrs.Name)

	if err == nil {
		return fmt.Errorf("package already installed bomm")

	}


	// ehck if pakge supports multiple version 


	ok, err := m.repo.PackageExists(agrs.Name)

	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("package not found")
	}

	


	pf ,err := m.repo.LoadPackage(agrs.Name,agrs.Version)
	if err != nil {
		return err
	}

	if pf.Release.Source.ARM64 == nil {
		fmt.Println("no ARM64 source found  ")
	}

	if pf.Release.Source.X64 == nil {
		fmt.Println("no X64 source found  ")
	}

	if pf.Release.Source.X86 == nil {
		fmt.Println("no X86 source found  ")
	}


	return nil
}
