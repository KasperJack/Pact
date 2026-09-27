package core

import (
	"fmt"
)



type Visitor interface {
	VisitShortcut(Shortcut) error
	VisitCommand(Command) error
	VisitAddPath(AddPath) error
}


type Block interface {
	BlockID() string
	//Run() error
	Accept(Visitor) error
	Name() string
	SuppotedScopes() []Scope
}









type Shortcut struct {
	ID string
	DisplayName string //optianl // // trim tarling and ending spacse 
	Exe         string // required // trim tarling and ending spacse  //check in path 
	Icon        string //optianl // trim tarling and ending spacse    //check in path
	Args        string //optianl // // trim tarling and ending spacse 
}

func (s Shortcut) Accept(v Visitor) error { return v.VisitShortcut(s) }




func (s Shortcut) Name() string {
	return "shortcut"
}

func (s Shortcut) SuppotedScopes() []Scope{
	return []Scope{ScopeUser,ScopeSystem}
}


func (s Shortcut) BlockID() string {
	return s.ID
}














type Command struct {
	ID   string
	Exe  string  // required // trim tarling and ending spacse  //check in path
	Args string //optianl // // trim tarling and ending spacse 
}


func (c Command) Name() string {
	return "command"
}

func (c Command) SuppotedScopes() []Scope{
	return []Scope{ScopeUser,ScopeSystem}
}

func (c Command) BlockID() string {
	return c.ID
}

func (c Command) Accept(v Visitor) error  { return v.VisitCommand(c) }















type AddPath struct {
	ID string
	Dir string // required // trim tarling and ending spacse  //check in path
}

func (a AddPath) Name() string {
	return "add_path"
}

func (a AddPath) SuppotedScopes() []Scope{
	return []Scope{ScopeSystem}
}

func (a AddPath) BlockID() string {
	return a.ID
}

func (a AddPath) Accept(v Visitor) error  { return v.VisitAddPath(a) }









type DryRunner struct{}

func (DryRunner) VisitShortcut(s Shortcut) error {
	fmt.Printf("[dry-run] would create shortcut %q -> %s\n", s.ID, s.Exe)
	return nil
}
func (DryRunner) VisitCommand(c Command) error {
	fmt.Printf("[dry-run] would run command: %s\n", c.Exe)
	return nil
}
func (DryRunner) VisitAddPath(a AddPath) error {
	fmt.Printf("[dry-run] would add to PATH: %s\n", a.Dir)
	return nil
}



type DryRemover struct{}
















type ManifestScope struct {
	InstallPath string
	Blocks      []Block // all types, file order preserved
}

type Manifest struct {
	Scope map[Scope]ManifestScope // key is scope name, e.g. "user" or "system"

}