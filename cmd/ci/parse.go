package main

import (
	//"github.com/hashicorp/hcl/v2/hclsimple"
    //"Pact/corelib/model"
    //"Pact/corelib/pipeline/publisher"
	//"Pact/corelib/client"
    //"github.com/yuin/gopher-lua"
    "Pact/corelib/luautil"

)

/*
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

        p.RebuildIndex()
        return nil




    default:
        return fmt.Errorf("pacakge not implmanted yet") 

    }


}
*/


func parseluaConfig(data []byte) error {

    /*
	s := client.NewLFilesystemSource("C:\\Users\\Aya\\Desktop\\pact-tools\\test-buckets\\defult")

	client := client.NewRepoClient(s)

	_ , err := publisher.NewPackage(data,client)

	if err != nil {
		return err
	}*/

	t := luautil.NewTestEvalContext()

     err := t.Eval(data)
    if err != nil {
        return err
    }

    err = t.RunInstall()

    if err != nil {
        return err
    }


	return nil
}



