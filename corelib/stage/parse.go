package stage

import (
	"fmt"
	//"github.com/BurntSushi/toml"
    //"slices"
    "github.com/hashicorp/hcl/v2/hclsimple"

)

type Processor interface {
    Validate() error
}

type Config struct {
  Release  *Release  `hcl:"release,block"`
  Package  *Package  `hcl:"package,block"`

}



func ParseConfig(raw []byte) (Processor, error) {
    var config Config

    err := hclsimple.Decode("package.hcl", raw, nil, &config)
    if err != nil {
        return nil, err
    }

    switch {
    case config.Package != nil && config.Release != nil:
        return nil, fmt.Errorf("you can only define a Package or a Release")
    case config.Package != nil:
        return config.Package, nil
    case config.Release != nil:
        return config.Release, nil
    default:
        return nil, fmt.Errorf("nothing was defined")
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
