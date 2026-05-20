package main


import (


	"os"
	"Pact/corelib/index"
	"fmt"

)

func main(){

	if len(os.Args) <= 2 {
		fmt.Println("expected command: build <arg>")
		os.Exit(1)
	}

	switch os.Args[1] {

	case "build":
		fmt.Printf("building %s \n","Package")
		build_cmd(os.Args[2])
		
	default:
		fmt.Println("expected command: build <arg>")
		os.Exit(1)

	}

}



func build_cmd (path string){

	index.LoadToml()
	
	//index.Ass()
	/*
	err := build(path)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
*/

	fmt.Println("ok")
	fmt.Println(path)

}