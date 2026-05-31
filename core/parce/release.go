package parce

import (
	
	"github.com/kasperjack/pact/core/model"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

func Release(rlsData []byte) (model.Release,error) {


	var config model.Release
	
    err := hclsimple.Decode("release.hcl", rlsData, nil, &config)
    if err != nil {
        return  model.Release{},err
    }

	return config,nil
}