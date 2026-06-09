package main

import (
	"fmt"
	"os"
	"runtime"

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

	
	
	
	var arch string

	if len(os.Args) > 4 {
		
		arch = os.Args[4]

		if err := validateArch(arch); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}


	}

	
	err := install(pkg,version,arch)

	if err != nil {
		fmt.Fprintln(os.Stderr,err)
		os.Exit(1)
	}


	fmt.Println("everything run ok")
	
}



func validateArch(target string) error {
    host := runtime.GOARCH

    switch host {
    case "arm64":
        if target != "arm64" {
            return fmt.Errorf("not supported install %s for an arm64 system", target)
        }
    case "amd64":
        if target != "x64" && target != "x86" {
            return fmt.Errorf("not supported install %s for a x64 system", target)
        }
    case "386":
        if target != "x86" {
            return fmt.Errorf("not supported install %s for a x86 system", target)
        }
    default:
        return fmt.Errorf("unsupported host architecture: %s", host)
    }

    return nil
}