package parce

import (
	"strings"
	"github.com/kasperjack/pact/core"
)


type ResourceBlock interface {
	export() core.Resource
}
	



var localResourceConstructors = map[string]func() ResourceBlock{

	core.Shortcut{}.Kind(): func() ResourceBlock { return &ShortcutBlock{}},
	core.Command{}.Kind(): func() ResourceBlock { return &CommandBlock{}},
	core.AddPath{}.Kind(): func() ResourceBlock { return &AddPathBlock{}},
	
}




///////////////////////////

type ShortcutBlock struct {
	Bound  *string `hcl:"bind,optional"`


	DisplayName *string `hcl:"display_name,optional"`
	Exe         string  `hcl:"exe"`
	Icon        *string `hcl:"icon,optional"`
	Args        *string `hcl:"args,optional"`


}

func (s *ShortcutBlock) export() core.Resource {

		return core.Shortcut{
		Exe:         strings.TrimSpace(s.Exe),
		DisplayName: strings.TrimSpace(derefOr(s.DisplayName, "")),
		Icon:        strings.TrimSpace(derefOr(s.Icon, "")),
		Args:        strings.TrimSpace(derefOr(s.Args, "")),

		Bound: core.Bound{
			Option: derefOr(s.Bound, ""),
		},
	}
}





////////////////////////////////


type CommandBlock struct {
	Bound  *string `hcl:"bind,optional"`


	Exe  string  `hcl:"exe"`
	Args *string `hcl:"args,optional"`

}


func (c *CommandBlock) export() core.Resource {


		return core.Command{

		Exe:  strings.TrimSpace(c.Exe),
		Args: strings.TrimSpace(derefOr(c.Args, "")),

		Bound: core.Bound{
			Option: derefOr(c.Bound, ""),
		},
	}
}






////////////////////


type AddPathBlock struct {

	Bound  *string `hcl:"bind,optional"`

	Dir string `hcl:"dir"`

}


func (a *AddPathBlock) export() core.Resource {

	
	return core.AddPath{
		Dir: strings.TrimSpace(a.Dir),
		Bound: core.Bound{
			Option: derefOr(a.Bound, ""),
		},
	}
}











