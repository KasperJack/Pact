package main

import (
	"fmt"
	"os"
	"path"
	//"path/filepath"

	"github.com/kasperjack/pact/core"
)

func install (pkg string, version string) error {

	bucket := "test-buckets"
	pkgPath := path.Join(bucket,pkg)

	info, err := os.Stat(pkgPath)
	if err != nil {
		return err
	}

	if  !info.IsDir() {
		return fmt.Errorf("path isnot a dir")
	}



	pkgFilePath :=  path.Join(pkgPath,fmt.Sprintf("%s.hcl",pkg))
	pkgFile, err := os.Open(pkgFilePath)
	if err != nil {
		return err
	}
	defer pkgFile.Close()



	releaseFilePath := path.Join(pkgPath,version,"release.hcl")
	releaseFile, err := os.Open(releaseFilePath)
	if err != nil {
		return err
	}
	defer releaseFile.Close()

	sciptFilePath := path.Join(pkgPath,version,"script.lua")
	scriptFile, err := os.Open(sciptFilePath)
	if err != nil {
		return err
	}
	defer scriptFile.Close()



	i := core.Input{
    Package: pkgFile,
    Release: releaseFile,
    Script:  scriptFile,
}
	

	
	

	return nil
}


// lib will take 3 manifest of the pcakge + install location that the pkmg uses 
// it will retrun an interface install update remove 