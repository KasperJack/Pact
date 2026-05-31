package model



type Package struct {
    Identifier  string `hcl:"identifier"`
    Name        string `hcl:"name"`
    Versioning  string `hcl:"versioning"`
    Description string `hcl:"description,optional"`
    Homepage    string `hcl:"homepage,optional"`
    License     string `hcl:"license,optional"`
}