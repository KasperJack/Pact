package core
import (
	"fmt"
	"github.com/kasperjack/pact/core/platform"
	//"github.com/kasperjack/pact/core/model"
	//"runtime"
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
	
	//var binSource string
	//var binHash string

	binSource,binHash,err := resolveSource(agrs.Arch,pf)
	if err != nil {
		return err
	}


	fmt.Println(binSource)
	fmt.Println(binHash)

	return nil
}












func resolveSource(arch platform.Arch, pf PackageFiles) (binSource string,binHash string,err error) {


	if arch == "" {
		fmt.Println("arch was not passed")
		host,err := platform.HostArch()

		if err != nil{ // this error sould not exist 
			return "","",err
		}



		switch host {
			case platform.X64:
				if  pf.Release.Source.X64 != nil {
					return pf.Release.Source.X64.URL,pf.Release.Source.X64.SHA256,nil
				}
				if  pf.Release.Source.X86 != nil {
					return pf.Release.Source.X86.URL,pf.Release.Source.X86.SHA256,nil
				}
				
				return "","",fmt.Errorf("no source found for x64/x86")

			case platform.X86:
				if  pf.Release.Source.X86 != nil {
					return pf.Release.Source.X86.URL,pf.Release.Source.X86.SHA256,nil
				}
				return "","",fmt.Errorf("no source found for x86")
			
			case platform.ARM64:
				if pf.Release.Source.ARM64 != nil {
					return pf.Release.Source.ARM64.URL,pf.Release.Source.ARM64.SHA256,nil
				}

				return "","",fmt.Errorf("no source found for arm64")

		}

	}else{
		fmt.Println("arch was passed")
		switch arch {
			case platform.X64:
				if  pf.Release.Source.X64 != nil {
					return pf.Release.Source.X64.URL,pf.Release.Source.X64.SHA256,nil
				}
				return "","",fmt.Errorf("no source found for x64")

			case platform.X86:
				if  pf.Release.Source.X86 != nil {
					return pf.Release.Source.X86.URL,pf.Release.Source.X86.SHA256,nil
				}
				return "","",fmt.Errorf("no source found for x86")
			
			case platform.ARM64:
				if pf.Release.Source.ARM64 != nil {
					return pf.Release.Source.ARM64.URL,pf.Release.Source.ARM64.SHA256,nil
				}
				return "","",fmt.Errorf("no source found for arm64")
		}
	}

	return "","",fmt.Errorf("unexpected error happend")
}