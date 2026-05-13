package model



type Release struct {
    PackageIdentifier string `hcl:"package_identifier"`
	Url    string `hcl:"url"`
	Version string `hcl:"version"`
    Hash string `hcl:"hash"`
}

