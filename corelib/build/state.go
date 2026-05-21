package build

import (

    //"fmt"
    //"Pact/corelib/client"
    "Pact/corelib/model"
	"path/filepath"
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"fmt"
	"strings"
)


func ReadPackages(rootDir string) ([]model.PackageContent, error) {
    var result []model.PackageContent

    //  (re, vu, we)
    prefixes, err := os.ReadDir(rootDir)
    if err != nil {
        return nil, fmt.Errorf("failed to read root dir: %w", err)
    }

    for _, prefix := range prefixes {
        if !prefix.IsDir() {
            continue
        }

        
        packages, err := os.ReadDir(filepath.Join(rootDir, prefix.Name()))
        if err != nil {
            return nil, fmt.Errorf("failed to read prefix dir %s: %w", prefix.Name(), err)
        }

        for _, pkg := range packages {
            if !pkg.IsDir() {
                continue
            }

			// sanity check:
			if !strings.HasPrefix(pkg.Name(), prefix.Name()) {
				return nil, fmt.Errorf("package %s is in wrong prefix folder %s", pkg.Name(), prefix.Name())
			}

            pkgPath := filepath.Join(rootDir, prefix.Name(), pkg.Name())
            pc, err := loadPackage(pkgPath)
            if err != nil {
                return nil, fmt.Errorf("failed to load package %s: %w", pkg.Name(), err)
            }

            result = append(result, pc)
        }
    }

    return result, nil
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

	if pkg.PackageIdentifier != filepath.Base(path) {
    	return model.PackageContent{}, fmt.Errorf(
        "package identifier %q does not match folder name %q", pkg.PackageIdentifier, filepath.Base(path),
    )
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
		 

		if r.Version != entry.Name() {
    		return model.PackageContent{}, fmt.Errorf(
        	"release version %q does not match folder name %q", r.Version, entry.Name(),
    	)
		}

        releases = append(releases, r)
    }



	return model.PackageContent{Package: pkg, Releases: releases}, nil
}