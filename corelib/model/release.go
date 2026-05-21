package model



type Release struct {
    PackageIdentifier string `hcl:"package_identifier"`
	Url    string `hcl:"url"`
	Version string `hcl:"version"`
    Hash string `hcl:"hash"`
}

func (r *Release) ToDomain() ReleaseT {

    rt := ReleaseT{
		Url: r.Url,
		Version: r.Version,
		Hash: r.Hash,

    }
    return rt
}

type ReleaseT struct {
	Url    string `hcl:"url" toml:"url"`
	Version string `hcl:"version" toml:"version"`
    Hash string `hcl:"hash" toml:"hash"`
}

func (r ReleaseT) ValidateOnRead () error {
	return nil
}


type ReleaseDB struct {
    PackageIdentifier string `gorm:"primaryKey;not null;references:package_identifier"`
    Version           string `gorm:"primaryKey"`
    Hash              string `gorm:"not null"`
    URL               string `gorm:"not null"`
}
func (ReleaseDB) TableName() string { return "releases" }