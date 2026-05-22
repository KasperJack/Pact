package model


type Package struct {
	PackageIdentifier string
	Name              string
	Versioning        string
	Description       string
	Homepage          string
	License           string
}


type PackageDB struct {
    PackageIdentifier string    `gorm:"primaryKey"`
    Name              string    `gorm:"not null"`
    Description       string
    Versioning        string    `gorm:"not null"`
    Releases          []ReleaseDB `gorm:"foreignKey:PackageIdentifier"`
}
func (PackageDB) TableName() string { return "packages" }
