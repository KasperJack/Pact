package main

import (
	"fmt"

	"os"

	"github.com/kasperjack/pact/core"
	"github.com/kasperjack/pact/core/parce"
)




func main(){

	_,err := verifyManifest("mfs.hcl")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}



	_,err = verifyIntface("inter.hcl")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}





	os.Exit(0)
}


func verifyIntface(path string)(*core.Interface, error) {


	i,err := parce.Interface(path)

	if err != nil {
		return nil,err
	}


	op := i.System.Options["4k_text"]

	fmt.Println(op.Description)
	return i,nil
}



func verifyManifest(path string)(*core.Manifest, error) {


	m,err := parce.Manifest(path)

	if err != nil {
		return nil,err
	}

	for _,s := range m.User.Shortcuts.Unconditional {

		fmt.Println(s.DisplayName)
	}



	return m,nil
}