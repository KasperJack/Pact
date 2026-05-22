package main

import (
	"fmt"
	//"github.com/hashicorp/hcl/v2/hclsimple"
    //"Pact/corelib/model"
    //"Pact/corelib/pipeline/publisher"
	//"Pact/corelib/client"
    "github.com/yuin/gopher-lua"

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
type Package struct {
	PackageIdentifier string
	Name              string
	Versioning        string
	Description       string
	Homepage          string
	License           string
}

func parseluaConfig(path string) error {
	L := lua.NewState()
	defer L.Close()

	var result Package

	// We define the "package" function in Lua
	L.SetGlobal("package", L.NewFunction(func(L *lua.LState) int {

		// first argument is the table passed to package { ... }
		tbl := L.CheckTable(1)

		result = Package{
			PackageIdentifier: tbl.RawGetString("package_identifier").String(),
			Name:              mustString(tbl, "name"),
			Versioning:        tbl.RawGetString("versioning").String(),
			Description:       tbl.RawGetString("description").String(),
			Homepage:          tbl.RawGetString("homepage").String(),
			License:           tbl.RawGetString("license").String(),
		}

		return 0
	}))

	// run lua file
	if err := L.DoFile(path); err != nil {
		panic(err)
	}

	// use result in Go
	fmt.Printf("%+v\n", result)

  

    return nil
}

func mustString(tbl *lua.LTable, key string) string {
    v := tbl.RawGetString(key)
    if v == lua.LNil {
        panic("missing required field: " + key)
    }
    return v.String()
}