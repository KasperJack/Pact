package core


import(
	//"fmt"
	"github.com/hashicorp/hcl/v2"


)



type ResolvedBlock interface {
	blockID() string
	blockRange() hcl.Range
	isResolvedBlock()
}

type Common struct {
	ID    string
	Range hcl.Range
}

func (c Common) blockID() string       { return c.ID }
func (c Common) blockRange() hcl.Range { return c.Range }
func (c Common) isResolvedBlock()      {}

type Shortcut struct {
	Common
	DisplayName string //optianl // // trim tarling and ending spacse 
	Exe         string // required // trim tarling and ending spacse  //check in path 
	Icon        string //optianl // trim tarling and ending spacse    //check in path
	Args        string //optianl // // trim tarling and ending spacse 
}

type Command struct {
	Common
	Exe  string  // required // trim tarling and ending spacse  //check in path
	Args string //optianl // // trim tarling and ending spacse 
}

type AddPath struct {
	Common
	Dir string // required // trim tarling and ending spacse  //check in path
}

type ResolvedScope struct {
	InstallPath string
	Blocks      []ResolvedBlock // all types, file order preserved
}

type ResolvedManifest struct {
	User   *ResolvedScope
	System *ResolvedScope
}