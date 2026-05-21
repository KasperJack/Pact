package model



type Release struct {
    PackageIdentifier string `hcl:"package_identifier"`
	Url    string `hcl:"url"`
	Version string `hcl:"version"`
    Hash string `hcl:"hash"`
}



type ReleaseT struct {
	Url    string `hcl:"url" toml:"url"`
	Version string `hcl:"version" toml:"version"`
    Hash string `hcl:"hash" toml:"hash"`
}
