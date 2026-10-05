package parce


import (

	"github.com/hashicorp/hcl/v2"
	"github.com/kasperjack/pact/core" )







type PackageBlock struct {

    Package    string  // check valid package identifier also matches with global package identifier
    Version  string `hcl:"version"` // check valid version
    Revision int `hcl:"revision"` // check valid revision
	Arch     string `hcl:"arch"` // check valid arch also matches with global pkg archs
	Scope    string `hcl:"scope"` // check valid scope also matches with global  pkg scopes

}


func (p *PackageBlock) SetLabels(labels []string) hcl.Diagnostics {
    if len(labels) != 1 {
        // diagnostic
    }

    p.Package = labels[0]
    return nil
}



func (p PackageBlock) export() core.Package {

	arch, _ := core.ParseArch(p.Arch)
	scope, _ := core.ParseScope(p.Scope)


	return core.Package{
		Identifier:    p.Package,
		Version:  p.Version,
		Revision: p.Revision,
		Arch:     arch,
		Scope:    scope,
	}
}



////////////////////////

type SourceBlock struct {

	Type   string 
	URL    string `hcl:"url"`	// check valid url	
	SHA256 string `hcl:"sha256"` // check valid sha256
	Size   int64  `hcl:"size"`

}
 

func (s *SourceBlock) SetLabels(labels []string) hcl.Diagnostics {
    if len(labels) != 1 {
        // diagnostic
    }

    s.Type = labels[0]
    return nil
}


func (s SourceBlock) export() core.Source {

	return core.Source{
		Type:   core.SourceZip,
		URL:    s.URL,
		SHA256: s.SHA256,
		Size:   s.Size,
	}
}








