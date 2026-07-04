package main

import (
	"fmt"
	"os"
	//"github.com/kasperjack/pact/core/platform"
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

	
	var arch string

	if len(os.Args) > 4 {
		arch = os.Args[4]
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








