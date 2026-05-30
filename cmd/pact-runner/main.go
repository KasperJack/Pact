package main


import (

	"os"
	"fmt"

)

func main(){

	if len(os.Args) <= 2 {
		fmt.Println("expected command: install <arg>")
		os.Exit(1)
	}

	switch os.Args[1] {

	case "build":
		fmt.Printf("install %s \n","Package")
		install_cmd(os.Args[2])
		
	default:
		fmt.Println("expected command: install <Package>")
		os.Exit(1)

	}

}



func install_cmd (pkg string){

	//index.LoadToml()
	
	//index.Ass()
	
	err := install(pkg)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}


	fmt.Println("ok")
	

}
