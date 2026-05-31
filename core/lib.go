package core


import (
	"io"
	//"fmt"
)


type PackageFiles struct {
    Package  io.Reader
    Release    io.Reader
    LuaScript  io.Reader
}




type LocalState interface {
	// this is an fs pov
	GetLockFile() (io.ReadWriter,error)
	CreatePackage(string) error
	GetPackagePath()
	PackageExists(string) bool

}


type Repo interface {

	PackageExists(string) (bool,error)
	LoadPackage(string,string) (PackageFiles,error)
	GetVersions(string) ([]string,error)

}