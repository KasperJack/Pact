package parce

import (
	
	"github.com/kasperjack/pact/core"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

func PackageInfo(pkgData []byte) (core.PackageInfo,error) {


	var config *core.PackageInfo = &core.PackageInfo{}
	
    err := hclsimple.Decode("package.hcl", pkgData, nil, config)
    if err != nil {
        return  core.PackageInfo{},err
    }


	err = config.Validate()

	if err != nil {
        return  core.PackageInfo{},err
    }



	return *config,nil
}