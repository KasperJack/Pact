package stage

import (
	"fmt"
	//"github.com/BurntSushi/toml"
    //"slices"
    "github.com/hashicorp/hcl/v2/hclsimple"

)

type Dependency struct {
	Name    string `hcl:"name,label"`
	Version string `hcl:"version"`
}

type Package struct {
	Name    string `hcl:"name"`
	Version string `hcl:"version"`
}

type Manifest struct {
	Package      Package      `hcl:"package,block"`
	Dependencies []Dependency `hcl:"dependency,block"`
}


func ParseConfig(raw []byte) (*Manifest, error) {
    var manifest Manifest



    err := hclsimple.Decode("package.hcl", raw, nil,&manifest)
	if err != nil {
		return nil,err
	}

	fmt.Println("Package:", manifest.Package.Name, manifest.Package.Version)

	for _, dep := range manifest.Dependencies {
		fmt.Printf("Dependency: %s @ %s\n", dep.Name, dep.Version)
	}






    return &manifest, nil
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
