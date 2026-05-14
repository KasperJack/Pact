package main

import (
	"fmt"
	"github.com/hashicorp/hcl/v2/hclsimple"
    "Pact/corelib/model"
    "Pact/corelib/pipeline/publisher"
	"Pact/corelib/client"

)

func parseConfig(raw []byte) (publisher.Pipeline, error) {

	s := client.NewLFilesystemSource("C:\\Users\\Aya\\Desktop\\pact-tools\\test-buckets\\defult")

	client := client.NewRepoClient(s)

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
        return publisher.NewPackage(config.Package,client),nil 

    default:
        return publisher.NewRelease(config.Release,client),nil 

    }


}