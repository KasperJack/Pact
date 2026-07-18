package manager

import (
	"fmt"
	//"path"
	"os"

	"path/filepath"

	//"github.com/kasperjack/pact/core/model"
	//"runtime"
	"github.com/kasperjack/pact/core/internal/runtime"
	//"github.com/nyaosorg/go-windows-junction"
	"github.com/kasperjack/pact/core/internal/win"
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
	resolvedRelease,resolvedArch,err := resolveSource(agrs.TargetArch,bundle.Release.Source)
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

	
	err = Download(resolvedRelease.URL,resolvedRelease.SHA256,stagingDirPath) // download cehck hash 
	if err != nil {
		return err
	}

	err = Extract(stagingDirPath) // extract and delete zip 
	if err != nil {
		return err
	}


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

	//// wtfffff
	currentDir := filepath.Join(filepath.Dir(installDir),"current")
	err = win.CreateJunction(installDir,currentDir)
	if err != nil {
		return err
	}

	// creating shortcuts
	if len(bundle.Package.Shortcuts) > 0 { 
		for _, sc := range bundle.Package.Shortcuts {
			exePath := filepath.Join(installDir,sc.Exe)
			// check if exe path exists 
			if _, err := os.Stat(exePath); err != nil{
				return fmt.Errorf("shortcut exe not found: %s", sc.Exe)
			}

		}


		for _, sc := range bundle.Package.Shortcuts {
			exePath := filepath.Join(currentDir,sc.Exe)

			var name string = sc.Name
			//var arguments string = sc.Args
			var iconPath string = sc.Icon

			if name == "" {
				name = bundle.Package.Name // should be fine for now Package Name is not optional
			}

			if iconPath == "" {
				iconPath = exePath
			}else{
				iconPath = filepath.Join(currentDir,sc.Icon)
			}



			win.CreateDesktopShortcut(name,exePath,sc.Args,iconPath)
		}

	}


	//test
	// passing a directory instead of a file path fails silently 
	win.CreateShim("C:\\Users\\Aya\\Desktop\\pact\\bin\\shims\\calc.exe","C:\\Windows\\System32\\calc.exe","")


	//TODO: commands exposed as shims
	if len(bundle.Package.Commands) > 0 { 

		for _, cmd := range bundle.Package.Commands {
			fmt.Println(cmd.Exe)
		}

	}
	/*
	//err = junction.Create()
	//fmt.Println(installDir)
	currentDir := filepath.Join(filepath.Dir(installDir),"current")

	//fmt.Println(currentDir)
	err = win.CreateJunction(installDir,currentDir) // C:\Users\installed\windirstat\2.6.1 , C:\Users\installed\windirstat\current

	if err != nil {
		return err
	}

	exePath := filepath.Join(currentDir,"WinDirStat.exe") // get pcakge exports for shortcts 


							//if exists name fall back to exe name or pcakge id 
	win.CreateDesktopShortcut(bundle.Package.Name,exePath,exePath)
	*/

	return nil
}











