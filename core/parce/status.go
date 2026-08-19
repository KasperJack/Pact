package parce

import (
	
	"github.com/kasperjack/pact/core"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

func ArchStatus(statusData []byte) (core.ArchStatus,error) {


	var config core.ArchStatus
	
    err := hclsimple.Decode("status.hcl", statusData, nil, &config)
    if err != nil {
        return  core.ArchStatus{},err
    }

	return config,nil
}