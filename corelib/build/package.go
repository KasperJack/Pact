package build
import (

    //"fmt"
    //"Pact/corelib/client"
    "Pact/corelib/model"
	//"path/filepath"
)



func Package(p *model.Package, bucketPath string) error {
/*
	//var bucketPath string = "C:\\Users\\Aya\\Desktop\\pact-tools\\test-buckets\\defult"

	packageDir := filepath.Join(bucketPath,p.PackageIdentifier[:2], p.PackageIdentifier) // needs to be created 

	packageFilePath := filepath.Join(packageDir, "package.toml")
	releaseFilePath := filepath.Join(packageDir,"releases",p.InitRelease.Version,"release.toml")

	pkg,release := p.ToDomain()

	err := writeToml(packageFilePath,pkg) //p
	if err != nil {
		return err
	}

	err = writeToml(releaseFilePath,release) //p.InitRelease
	if err != nil {
		return err
	}

*/
	return nil
}