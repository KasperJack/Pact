package stage


type Package struct {
    Name string `hcl:"name"`
    PackageIdentifier string `hcl:"package_identifier"`
    Description string `hcl:"description"`
    InitRelease *InitRelease `hcl:"release,block"`
}

type InitRelease struct {
	Url    string `hcl:"url"`
	Version string `hcl:"version"`
    Hash string `hcl:"hash"`
}

func (r *Package) Validate () error {
    //fmt.Println("validating package")
    return nil
}  
