package model


type Package struct {
	PackageIdentifier string `lua:"package_identifier"`
	Name              string `lua:"name"`
	Versioning        string `lua:"versioning"`
	Description       string `lua:"description"`
	Homepage          string `lua:"homepage"`
	License           string `lua:"license"`
}


