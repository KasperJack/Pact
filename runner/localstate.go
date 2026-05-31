package main

import(

	"github.com/kasperjack/pact/core"
	"io"
	"bytes"
)

type localState struct {
	userPackageRoot string
}



func NewLocalState(userPackageRoot string) core.LocalState {
 return &localState{
	userPackageRoot: userPackageRoot,
 }
}



func (l *localState) GetLockFile() (io.ReadWriter, error) {
    return &bytes.Buffer{}, nil
}


func (l *localState) PackageExists(packageName string) (bool, error) {

    return true, nil
}




