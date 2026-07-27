package main

import (
	"fmt"
	"os"
	"strings"
	"github.com/kasperjack/pact/core"
	//"golang.org/x/text/cases"
	//"github.com/kasperjack/pact/psbridge"
)








func main(){
	
	if len(os.Args) < 4 {
		fmt.Println("expected command: install <pkg> <ver>")
		os.Exit(1)
	}

	switch os.Args[1] {

	case "install":
		//fmt.Printf("installing %s \n", os.Args[2])
		installCmd(os.Args[2],os.Args[3])
		
	default:
		fmt.Println("expected command: install <pkg> <ver>")
		os.Exit(1)

	}

}



func installCmd (pkg string, version string){

	
	var arch core.Arch

	if len(os.Args) > 4 {
		var err error
		archflag := os.Args[4]

		arch, err = ParseArchFlag(archflag)
		if err != nil {
			fmt.Fprintln(os.Stderr,err)
			os.Exit(1)
		}

		if err := ValidateArchForHost(arch, core.HostArch()); err != nil {
			fmt.Fprintln(os.Stderr,err)
			os.Exit(1)
		}
		
	}

	//fmt.Printf("%s\n",arch)
	
	err := install(pkg,version,arch)

	if err != nil {
		fmt.Fprintln(os.Stderr,err)
		os.Exit(1)
	}

	//a := platform.Arch("stringFromCli")
	fmt.Println("everything run ok")
	
}




func ParseArchFlag(s string) (core.Arch, error) {
    switch strings.ToLower(strings.TrimSpace(s)) {
    case "":
        return core.ArchUndefined, nil
    case "x86", "32", "32bit", "386":
        return core.ArchX86, nil
    case "x64", "64", "64bit", "amd64":
        return core.ArchX64, nil
    case "arm64", "aarch64":
        return core.ArchArm64, nil
    default:
        return core.ArchUndefined, fmt.Errorf("unrecognized architecture %q", s)
    }
}





func ValidateArchForHost(requested, host core.Arch) error {
    if requested == core.ArchUndefined {
        return nil 
    }
	
    if requested == host {
        return nil
    }
   
    if requested == core.ArchX86 && host == core.ArchX64{
        return nil
    }

    return fmt.Errorf("architecture %s is not supported on this host (%s) — try omitting --arch to auto-select, or use a compatible architecture",requested.String(),host.String())
}







/*
func validateArch(targetArch platform.Arch) error {
    
	host, err := platform.HostArch()

	if err != nil {
		return err 
	}
	
    switch host {

		case platform.ARM64:
			if targetArch != platform.ARM64 {
				return fmt.Errorf("cannot install %s binary on an arm64 host", targetArch)
			}
	


		case platform.X64:
			if targetArch != platform.X64 && targetArch != platform.X86 {
				return fmt.Errorf("cannot install %s binary on an x64 host", targetArch)
			}
	

		case platform.X86:
			if targetArch != platform.X86 {
				return fmt.Errorf("cannot install %s binary on an x86 host", targetArch)
			}
		default:
    		return fmt.Errorf("unsupported host architecture: %s", host)
	}

	return nil
}*/








