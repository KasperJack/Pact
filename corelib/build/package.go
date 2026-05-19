package build
import (

    //"fmt"
    //"Pact/corelib/client"
    "Pact/corelib/model"
	"path/filepath"
)



func Package(p *model.Package) error {

	var bucketPath string = "C:\\Users\\Aya\\Desktop\\pact-tools\\test-buckets\\defult"

	packageDir := filepath.Join(bucketPath,p.PackageIdentifier[:2], p.PackageIdentifier) // needs to be created 

	packageFilePath := filepath.Join(packageDir, "package.toml")
	releaseFilePath := filepath.Join(packageDir,"releases",p.InitRelease.Version,"release.toml")

	err := writeToml(packageFilePath,p.InitRelease)
}