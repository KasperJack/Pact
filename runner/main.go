package main

import (
	"fmt"
	"os"
)

func main(){

	if len(os.Args) <= 3 {
		fmt.Println("expected command: install <pkg> <ver>")
		os.Exit(1)
	}

	switch os.Args[1] {

	case "install":
		fmt.Printf("install %s \n","Package")
		install_cmd(os.Args[2],os.Args[3])
		
	default:
		fmt.Println("expected command: install <pkg> <ver>")
		os.Exit(1)

	}

}



func install_cmd (pkg string, version string){

	//index.LoadToml()
	
	//index.Ass()
	
	err := install(pkg,version)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}


	fmt.Println("ok")
	

}
