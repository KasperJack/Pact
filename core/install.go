package core
import (
	"fmt"
	//"github.com/kasperjack/pact/core/model"
	"runtime"
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

	


	pf ,err := m.repo.LoadPackage(agrs.Name,agrs.Version)
	if err != nil {
		return err
	}
	
	var binSource string
	var binHash string

	//TODO: create a ResolveSource func
	if agrs.Arch == "" {
		host := runtime.GOARCH
		switch host {
			case "arm64":
				if pf.Release.Source.ARM64 == nil {
					return fmt.Errorf("no ARM64 source found")
				}
				binSource = pf.Release.Source.ARM64.URL
				binHash = pf.Release.Source.ARM64.SHA256


			case "amd64":
				if pf.Release.Source.X64 != nil {
					binSource = pf.Release.Source.X64.URL
					binHash = pf.Release.Source.X64.SHA256
				}else if pf.Release.Source.X86 != nil {
					binSource = pf.Release.Source.X86.URL
					binHash = pf.Release.Source.X86.SHA256
				}else {
					return fmt.Errorf("no x86/x64 source found")
				}
				

			case "386":
				if pf.Release.Source.X86 == nil {
					return fmt.Errorf("no x86 source found")
				}
				binSource = pf.Release.Source.X86.URL
				binHash = pf.Release.Source.X86.SHA256

			default:
				return fmt.Errorf("unsupported host architecture: %s", host)
				
    	}

	}else{

		fmt.Printf("looking for source %s \n",agrs.Arch)
	}



	fmt.Println(binSource)
	fmt.Println(binHash)

	return nil
}
