package main

import (
	"fmt"
	"github.com/hashicorp/hcl/v2/hclsimple"
    "Pact/corelib/model"
    "Pact/corelib/pipeline/verifier"

)

func parseConfig(raw []byte) (verifier.Pipeline, error) {
    var config model.Config

    err := hclsimple.Decode("package.hcl", raw, nil, &config)
    if err != nil {
        return nil, err
    }


    switch {
    case config.Package == nil && config.Release == nil:
        return nil,fmt.Errorf("mssing def")

    case config.Package != nil && config.Release != nil:
        return nil,fmt.Errorf("can't have a package and a release def at the same time")

    case config.Package != nil:
        return verifier.NewPackage(config.Package),nil 

    default:
        return verifier.NewRelease(config.Release),nil 

    }


}