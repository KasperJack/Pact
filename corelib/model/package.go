package model
import "fmt"

type Package struct {
    Name              string      `hcl:"name" toml:"name" gorm:"not null"`
    PackageIdentifier string      `hcl:"package_identifier" toml:"package_identifier" gorm:"primaryKey"`
    Description       string      `hcl:"description,optional" toml:"description,omitempty"`
    Versioning        string      `hcl:"versioning" toml:"versioning" gorm:"not null"`
    InitRelease       ReleaseT     `hcl:"release,block" toml:"-" gorm:"-"`
}






type PackageT struct {
    Name              string     `toml:"name"`
    PackageIdentifier string     `toml:"package_identifier"`
    Description       string     `toml:"description,omitempty"`
    Versioning        string     `toml:"versioning"`

}












func (p *Package) ValidateOnRead() error {
    	
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

