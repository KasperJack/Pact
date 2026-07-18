package core


import (
	//"github.com/kasperjack/pact/core/model"

	
)





/*
app  →  manager, core
manager  →  runtime, core
runtime  →  core
core  →  (leaf)

core(types and install args)
manager(all manager commands)








reorganize packages into core and manager to fix import cycles
- Add manager package to hold manager operations
- Move types from model into core
- Move manager-related logic out of core into manager
- Resolve import cycle between core and runtime






*/



type PackageBundle struct {
    Package  Package
    Release    Release
    Script  []byte
}


type LocalState interface {
	// this is an fs pov
	//CreatePackage(string) error
	CreateInstallDir(PackageIdentifier string, version string) (string,error)
	PackageExists(string) (bool, error)

}
type Repo interface {

	PackageExists(string) (bool,error)
	LoadPackage(PackageIdentifier string, version string) (PackageBundle,error)
	GetVersions(string) ([]string,error)

}

type LockFile interface {
    
    GetInstalled(PackageIdentifier string) (LockedPackage, error)
    RecordInstall(pkg LockedPackage) error
    RecordRemove(PackageIdentifier string, version string) error
	Test() error

}
    //InstallDir(pkg string) string
	//IsInstalled(pkg string) error







////////



type LockedPackage struct {
    Name        string `hcl:"name,label"`
    Version     string `hcl:"version"`
    InstalledAt string `hcl:"installed_at"`
    InstallDir  string `hcl:"install_dir"`
}

type LockFileC struct {
    Packages []LockedPackage `hcl:"package,block"`
}










type Shortcut struct {
    Name string `hcl:"name,optional"`
    Exe  string `hcl:"exe"`
    Icon string `hcl:"icon,optional"`
    Args string `hcl:"args,optional"`
}



type Command struct {
    Exe  string `hcl:"exe"`
    Args string `hcl:"args,optional"`   // default args baked into shim
}




type Package struct {   // size=144 (0x90) /use a pointer ? 
    Identifier  string `hcl:"identifier"`
    Name        string `hcl:"name"`
    Versioning  string `hcl:"versioning"`
    Description string `hcl:"description,optional"`
    Homepage    string `hcl:"homepage,optional"`
    License     string `hcl:"license,optional"`
    Shortcuts   []Shortcut `hcl:"shortcut,block"`
    Commands    []Command  `hcl:"command,block"`
}














type ReleaseSource struct {
    URL    string `hcl:"url"`
    SHA256 string `hcl:"sha256"`
}

type ReleaseSourceBlock struct {
    X64      *ReleaseSource `hcl:"x64,block"`
    ARM64    *ReleaseSource `hcl:"arm64,block"`
    X86      *ReleaseSource `hcl:"x86,block"`
    Universal *ReleaseSource `hcl:"universal,block"`
    NoArch   *ReleaseSource `hcl:"noarch,block"`
}

type Release struct {
    Version string             `hcl:"version"`
    Source  ReleaseSourceBlock `hcl:"source,block"`
}