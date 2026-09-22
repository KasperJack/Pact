package core

import (
	"fmt"
)

type Block interface {
	BlockID() string
	Run() error
	Name() string
	SuppotedScopes() [2]Scope
}






type Shortcut struct {
	ID string
	DisplayName string //optianl // // trim tarling and ending spacse 
	Exe         string // required // trim tarling and ending spacse  //check in path 
	Icon        string //optianl // trim tarling and ending spacse    //check in path
	Args        string //optianl // // trim tarling and ending spacse 
}


func (s Shortcut) Run() error {
	return nil
}

func (s Shortcut) Name() string {
	return "shortcut"
}

func (s Shortcut) SuppotedScopes()[2]Scope{
	return [2]Scope{ScopeUser,ScopeSystem}
}


func (s Shortcut) BlockID() string {
	return s.ID
}








type Command struct {
	ID   string
	Exe  string  // required // trim tarling and ending spacse  //check in path
	Args string //optianl // // trim tarling and ending spacse 
}

func (c Command) Run() error {
	return nil
}

func (c Command) Name() string {
	return "command"
}

func (c Command) SuppotedScopes()[2]Scope{
	return [2]Scope{ScopeUser,ScopeSystem}
}

func (c Command) BlockID() string {
	return c.ID
}





type AddPath struct {
	ID string
	Dir string // required // trim tarling and ending spacse  //check in path
}

func (a AddPath) Run() error {

	fmt.Println("adding path: ", a.Dir)
	return nil
}

func (a AddPath) Name() string {
	return "add_path"
}

func (a AddPath) SuppotedScopes()[2]Scope{
	return [2]Scope{ScopeUser,ScopeSystem}
}

func (a AddPath) BlockID() string {
	return a.ID
}









type ManifestScope struct {
	InstallPath string
	Blocks      []Block // all types, file order preserved
}

type Manifest struct {
	Scope map[Scope]ManifestScope // key is scope name, e.g. "user" or "system"

}