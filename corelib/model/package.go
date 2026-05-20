package model
import "fmt"

type Package struct {
    Name              string      `hcl:"name" toml:"name" gorm:"not null"`
    PackageIdentifier string      `hcl:"package_identifier" toml:"package_identifier" gorm:"primaryKey"`
    Description       string      `hcl:"description,optional" toml:"description,omitempty"`
    Versioning        string     `hcl:"Versioning" toml:"versioning" gorm:"not null"`
    InitRelease       InitRelease `hcl:"release,block" toml:"-" gorm:"-"`
}


func (p *Package) ValidateTomlRead() error {
    	
    if p.Name == "" {
		return fmt.Errorf("name is required")
	}

	if p.PackageIdentifier == "" {
		return fmt.Errorf("PackageIdentifier is required")
	}

    if p.Versioning == "" {
        return fmt.Errorf("Versioning is required")
    }

	return nil
}

type InitRelease struct {
    Url     string `hcl:"url" toml:"url"`
    Version string `hcl:"version" toml:"version"`
    Hash    string `hcl:"hash" toml:"hash"`
}

