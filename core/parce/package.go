package parce

import (
	
	"github.com/kasperjack/pact/core/model"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

func Pacakge(pkgData []byte) (model.Package,error) {


	var config model.Package
	
    err := hclsimple.Decode("package.hcl", pkgData, nil, &config)
    if err != nil {
        return  model.Package{},err
    }

	return config,nil
}