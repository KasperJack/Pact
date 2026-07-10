package manager

import (
	"fmt"



	//"github.com/kasperjack/pact/core/model"
	//"runtime"
	"github.com/kasperjack/pact/core/internal/runtime"
)




func (m *pkgManager) Install(agrs InstallArgs) error {

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

	


	bundle ,err := m.repo.LoadPackage(agrs.Name,agrs.Version)
	if err != nil {
		// error loading the package 
		return err
	}
	
	

	//fix this
	s,a,err := resolveSource(agrs.TargetArch,bundle.Release.Source)
	if err != nil {
		return err
	}


	fmt.Println(a)
	fmt.Println(s.URL)
	fmt.Println(s.SHA256)
	

	/*/////////////
	 //prepare staging dir
	stagingDirPath,err := m.staging.Prepare(agrs.Name,agrs.Version) // --> C:\pact\staging\ripgrep\14.1.0\
	defer m.staging.Clear()
	if err != nil {
		return err

	}

	err = Download(binSource,binHash,stagingDirPath) // download cehck hash 
	if err != nil {
		return err
	}

	err = Extract(stagingDirPath) // extract and delete zip 
		if err != nil {
		return err
	}*/



	// stagingDir , resolvedArch,PackageBundle, interfaces (create a lock entrey,localState a pointer to install dir)
	rt, err := runtime.NewIinstallContext(bundle)
	if err != nil {
		return err

	}

	err = rt.Run()

	if err != nil {
		return err
	}

	return nil
}











