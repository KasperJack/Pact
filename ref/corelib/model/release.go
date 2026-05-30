package model


type Release struct {
    PackageIdentifier string `hcl:"package_identifier"`
	Url    string `hcl:"url"`
	Version string `hcl:"version"`
    Hash string `hcl:"hash"`
}


type ReleaseDB struct {
    PackageIdentifier string `gorm:"primaryKey"`
    Version           string `gorm:"primaryKey"`
    Hash              string `gorm:"not null"`
    URL               string `gorm:"not null"`
}
func (ReleaseDB) TableName() string { return "releases" }