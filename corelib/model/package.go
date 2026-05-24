package model
import ("github.com/yuin/gopher-lua")

type Package struct {
	PackageIdentifier string `lua:"package_identifier"`
	Name              string `lua:"name"`
	Versioning        string `lua:"versioning"`
	Description       string `lua:"description"`
	Homepage          string `lua:"homepage"`
	License           string `lua:"license"`
	InstallFn         *lua.LFunction `lua:"install"`
}


type PackageDB struct {
    PackageIdentifier string    `gorm:"primaryKey"`
    Name              string    `gorm:"not null"`
    Description       string
    Versioning        string    `gorm:"not null"`
    Releases          []ReleaseDB `gorm:"foreignKey:PackageIdentifier"`
}
func (PackageDB) TableName() string { return "packages" }
