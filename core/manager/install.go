package manager

import (
	"fmt"



	//"github.com/kasperjack/pact/core/model"
	//"runtime"
	"github.com/kasperjack/pact/core/internal/runtime"
)




func (m *pkgManager) Install(agrs InstallArgs) error {

	_ , err := m.lockFile.GetInstalled(agrs.PackageIdentifier)

	if err == nil {
		return fmt.Errorf("package already installed boom")

	}


	//ehck if pakge supports multiple version 


	ok, err := m.repo.PackageExists(agrs.PackageIdentifier)

	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("package not found")
	}

	


	bundle ,err := m.repo.LoadPackage(agrs.PackageIdentifier,agrs.Version)
	if err != nil {
		// error loading the package 
		return err
	}
	
	

	//fix this
	_,resolvedArch,err := resolveSource(agrs.TargetArch,bundle.Release.Source)
	if err != nil {
		return err
	}


	//fmt.Println(a)
	//fmt.Println(s.URL)
	//fmt.Println(s.SHA256)
	

	
	 //prepare staging dir
	stagingDirPath,err := m.staging.Prepare(agrs.PackageIdentifier,agrs.Version) // --> C:\pact\staging\ripgrep\14.1.0\
	
	if err != nil {
		return err

	}
	defer m.staging.Clear()


	/*
	err = Download(resolvedRelease.URL,resolvedRelease.SHA256,stagingDirPath) // download cehck hash 
	if err != nil {
		return err
	}

	err = Extract(stagingDirPath) // extract and delete zip 
	if err != nil {
		return err
	}*/


	installDir, err := m.localState.CreateInstallDir(agrs.PackageIdentifier,agrs.Version)
	if err != nil {
		return err
	}
	// does pcakges an install dir(is portable)	
	// installDir ,stagingDir , resolvedArch,PackageBundle, interfaces (create a lock entrey,localState a pointer to install dir)
	rt, err := runtime.NewIinstallContext(installDir,stagingDirPath,resolvedArch,bundle)
	if err != nil {
		return err

	}

	err = rt.Run()

	if err != nil {
		return err
	}

	return nil
}











