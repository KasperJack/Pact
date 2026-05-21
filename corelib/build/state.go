package build

import (

    //"fmt"
    //"Pact/corelib/client"
    "Pact/corelib/model"
	"path/filepath"
	"github.com/BurntSushi/toml"
	"log"
	"os"
)


func buildState (rootDir string) ([]model.PackageContent,error) {

	return nil,nil
}






func loadPackage(path string) (model.PackageContent,error){
	

	var pkg model.PackageT

	packageFilePath := filepath.Join(path, "package.toml")

	_, err := toml.DecodeFile(packageFilePath, &pkg)
	if err != nil {
		log.Fatal(err)
	}

	err = pkg.ValidateOnRead()
	if err != nil {
		log.Fatal(err)
	}



	entries, err := os.ReadDir(filepath.Join(path,"releases"))
	if err != nil {
		return model.PackageContent{},err
	}

	var releases []model.ReleaseT

	for _, entry := range entries {
		
        if !entry.IsDir() {
            continue
        }

		var r model.ReleaseT
		_, err := toml.DecodeFile(filepath.Join(path,"releases",entry.Name(),"release.toml"), &r)
		if err != nil {
			log.Fatal(err)
		}

		 err = r.ValidateOnRead()
		 if err != nil {
			log.Fatal(err)
		 }


        releases = append(releases, r)
    }



	return model.PackageContent{Package: pkg, Releases: releases}, nil
}