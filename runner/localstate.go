package main

import(

	"github.com/kasperjack/pact/core"

)

type localState struct {
	baseDir string
}



func NewLocalState(baseDir string) core.LocalState {
 return &localState{
	baseDir: baseDir,
 }
}




func (l *localState) PackageExists(packageName string) (bool, error) {

    return true, nil
}




