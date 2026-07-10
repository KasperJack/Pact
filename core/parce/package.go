package parce

import (
	
	"github.com/kasperjack/pact/core"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

func Pacakge(pkgData []byte) (core.Package,error) {


	var config core.Package
	
    err := hclsimple.Decode("package.hcl", pkgData, nil, &config)
    if err != nil {
        return  core.Package{},err
    }

	return config,nil
}