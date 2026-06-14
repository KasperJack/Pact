package main

import (
	"fmt"
	"os"
	"github.com/kasperjack/pact/core/platform"
	//"golang.org/x/text/cases"
)








func main(){

	if len(os.Args) < 4 {
		fmt.Println("expected command: install <pkg> <ver>")
		os.Exit(1)
	}

	switch os.Args[1] {

	case "install":
		//fmt.Printf("installing %s \n", os.Args[2])
		install_cmd(os.Args[2],os.Args[3])
		
	default:
		fmt.Println("expected command: install <pkg> <ver>")
		os.Exit(1)

	}

}



func install_cmd (pkg string, version string){

	
	
	
	var arch platform.Arch

	if len(os.Args) > 4 {
		
		arch,err := parseArch(os.Args[4])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}


		if err := validateArch(arch); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}


	}

	fmt.Println(arch)

	/*
	err := install(pkg,version,arch)

	if err != nil {
		fmt.Fprintln(os.Stderr,err)
		os.Exit(1)
	}

	//a := platform.Arch("stringFromCli")
	fmt.Println("everything run ok")
	*/
}



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
}








func parseArch(s string) (platform.Arch, error) {
    switch platform.Arch(s) {
    case platform.X86, platform.X64, platform.ARM64:
        return platform.Arch(s), nil
    default:
        return "", fmt.Errorf("unknown arch %q, must be one of: x86, x64, arm64", s)
    }
}