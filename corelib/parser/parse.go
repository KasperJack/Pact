package parser

import (
	"fmt"
	"github.com/hashicorp/hcl/v2/hclsimple"
    "Pact/corelib/model"
    "Pact/corelib/pipeline"

)



func ParseConfig(raw []byte) (pipeline.Pipeline, error) {
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
        return pipeline.NewPackage(config.Package),nil 

    default:
        return pipeline.NewRelease(config.Release),nil 

    }


}

func ParseConfigCheck(raw []byte) (pipeline.Checker, error) {
    var config model.Config

    err := hclsimple.Decode("package.hcl", raw, nil, &config)
    if err != nil {
        return nil, err
    }

    // checks

        switch {
    case config.Package == nil && config.Release == nil:
        return nil,fmt.Errorf("mssing def")

    case config.Package != nil && config.Release != nil:
        return nil,fmt.Errorf("can't have a package and a release def at the same time")

    case config.Package != nil:
        return pipeline.NewPackage(config.Package),nil 

    default:
        return pipeline.NewRelease(config.Release),nil 

    }

    
}










    /*
    var root map[string]any
    if err := toml.Unmarshal(raw, &root); err != nil {
        return nil, err
    }
    root_keys := []string{"release","package"}


    for key := range root {
        if !slices.Contains(root_keys, key) {
            return  nil, fmt.Errorf("config: unknown key %s", key)
        }
    }


    _, ok := root["package"]
    // pacakge def 
    if ok {

        
        _, ok = root["release"]
    
        if ok {

            
        }
        
        // error 

    }


    _, ok = root["release"]
    
    if ok {

        fmt.Println("ok")
    }
  
    // no pacakge ,no release
    

    //fmt.Printf("%T\n", r)
    fmt.Printf("%T\n", p)





    /*

    var cfg Config
    if err := toml.Unmarshal(raw, &cfg); err != nil {
        return nil, err
    }


    switch {
    case cfg.Release != nil:
        // handle release
    case cfg.Package != nil:
        // handle package
    default:
        return nil, fmt.Errorf("config: must define either [release] or [package]")
    }*/
