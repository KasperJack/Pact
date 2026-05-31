package core


import (
	"io"
	"github.com/kasperjack/pact/core/model"
	"github.com/kasperjack/pact/core/parce"
)



func ParcePacakge(pkgData []byte) (model.Package,error) {

	return parce.Pacakge(pkgData)

}

func ParceRelease(rlsData []byte) (model.Release,error) {

	return parce.Release(rlsData)
}






type PackageFiles struct {
    Package  model.Package
    Release    model.Release
    LuaScript  []byte
}




type LocalState interface {
	// this is an fs pov
	GetLockFile() (io.ReadWriter,error)
	//CreatePackage(string) error
	//GetPackagePath()
	PackageExists(string) (bool, error)

}


type Repo interface {

	PackageExists(string) (bool,error)
	LoadPackage(packageName string, version string) (PackageFiles,error)
	GetVersions(string) ([]string,error)

}