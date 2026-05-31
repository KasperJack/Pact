package core


import (
	"io"
	"fmt"
)


type Input struct {
    PackageFile  io.Reader
    RelaseFile    io.Reader
    LuaScript  io.Reader
}

func Process(input Input) error {


	config, _ := io.ReadAll(input.PackageFile)
    data, _   := io.ReadAll(input.RelaseFile)
    schema, _ := io.ReadAll(input.LuaScript)

	fmt.Println(string(config))
	fmt.Println(string(schema))
	fmt.Println(string(data))
	return nil
}



type LocalState interface {

	GetLockFile() []byte
	LoadPackage()

}