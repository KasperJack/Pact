package core


import (
	"github.com/kasperjack/pact/core/model"

	
)



type PackageBundle struct {
    Package  model.Package
    Release    model.Release
    Script  []byte
}


type LocalState interface {
	// this is an fs pov
	//CreatePackage(string) error
	//GetPackagePath()
	PackageExists(string) (bool, error)

}
type Repo interface {

	PackageExists(string) (bool,error)
	LoadPackage(packageName string, version string) (PackageBundle,error)
	GetVersions(string) ([]string,error)

}

type LockFile interface {
    
    GetInstalled(pkg string) (model.LockedPackage, error)
    RecordInstall(pkg model.LockedPackage) error
    RecordRemove(pkg string) error
	Test() error

}
    //InstallDir(pkg string) string
	//IsInstalled(pkg string) error







