package main

import (
	"fmt"
	"github.com/hashicorp/hcl/v2/hclsimple"
    "Pact/corelib/model"
    "Pact/corelib/pipeline/publisher"
	"Pact/corelib/client"

)

func parseConfig(raw []byte) error {

	s := client.NewLFilesystemSource("C:\\Users\\Aya\\Desktop\\pact-tools\\test-buckets\\defult")

	client := client.NewRepoClient(s)

    var config model.Config

    err := hclsimple.Decode("package.hcl", raw, nil, &config)
    if err != nil {
        return  err
    }


    switch {
    case config.Package == nil && config.Release == nil:
        return fmt.Errorf("mssing def")

    case config.Package != nil && config.Release != nil:
        return fmt.Errorf("can't have a package and a release def at the same time")

    case config.Package != nil:
        p := publisher.NewPackage(config.Package,client)

        err := p.Validate()

        if err != nil {
            return  err
        }

        err = p.Build()

        if err != nil {
            return  err
        }
        return nil




    default:
        return fmt.Errorf("pacakge not implmanted yet") 

    }


}