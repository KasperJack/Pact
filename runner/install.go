package main

import (
	"os"
	"fmt"
	"reflect"
	
	"github.com/kasperjack/pact/core"
	"github.com/kasperjack/pact/core/manager"

)





func install(pkg string, version string, arch core.Arch) error {

	localState := NewLocalState()


	printStruct(localState)
	os.Exit(0)



	m,err := manager.NewManager(localState)
	if err != nil {
		return err
	}



	
	err = m.Install(core.InstallArgs{
		PackageIdentifier: pkg,
		Version:           core.ParseVersion(version),
		TargetArch:        arch,
	})

	
	if err != nil {
		return err
	}

	return nil
}


func printStruct(v any) {
    rv := reflect.ValueOf(v)
    if rv.Kind() == reflect.Pointer {
        rv = rv.Elem()
    }

    rt := rv.Type()

    for i := 0; i < rv.NumField(); i++ {
        fmt.Printf("%s: %v\n", rt.Field(i).Name, rv.Field(i).Interface())
    }
}