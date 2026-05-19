package model


type Package struct {
    Name              string      `hcl:"name" toml:"name"`
    PackageIdentifier string      `hcl:"package_identifier" toml:"package_identifier"`
    Description       string      `hcl:"description,optional" toml:"description,omitempty"`
    Versioning        string     `hcl:"Versioning" toml:"versioning"`
    InitRelease       InitRelease `hcl:"release,block" toml:"-"`
}

type InitRelease struct {
    Url     string `hcl:"url" toml:"url"`
    Version string `hcl:"version" toml:"version"`
    Hash    string `hcl:"hash" toml:"hash"`
}

