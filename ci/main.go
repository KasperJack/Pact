package main

import (
	"fmt"
	"os"
	//"path/filepath"

	//"github.com/knqyf263/go-deb-version"
)


func main(){


	if len(os.Args) > 2 {
		
		switch os.Args[1] {

		case "build":

			buildCmd(os.Args[2])

			
		}
	
	}


	fmt.Fprintln(os.Stderr,"use build <pkg>")
	os.Exit(1)
}



func buildCmd(pkg string){


}