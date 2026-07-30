package parce

import (
	
	"github.com/kasperjack/pact/core"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

func Release(rlsData []byte) (core.Release,error) {


	var config core.Release
	
    err := hclsimple.Decode("package.hcl", rlsData, nil, &config)
    if err != nil {
        return  core.Release{},err
    }

	return config,nil
}