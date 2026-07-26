package core

import (
	//"github.com/kasperjack/pact/core/model"

	"errors"

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




type InstallArgs struct {
    PackageIdentifier string
    Version    Version
	TargetArch  Arch
}







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

	PackageExists(PackageIdentifier string) (bool,string,error)
	LoadPackage(PackageIdentifier string, version string) (PackageBundle,error)
	GetVersions(PackageIdentifier string) ([]string,error)
    GetLatest(PackageIdentifier string) (string, error)

}

type LockFile interface {
    
    GetInstalled(PackageIdentifier string) (map[string]LockedPackage)
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
    Arch        string `hcl:"arch"`
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






////version







type versionKind int

const (
	versionUndefined versionKind = iota
	versionExact
)

// Version represents a package version, which may be undefined
// (caller didn't specify one — resolve to latest downstream).
type Version struct {
	kind  versionKind
	value string
}

func versionUndefinedValue() Version {
	return Version{kind: versionUndefined}
}

func versionExactValue(v string) Version {
	if v == "" {
		panic("core.versionExactValue: empty version string")
	}
	return Version{kind: versionExact, value: v}
}

// ParseVersion converts a raw string (from a CLI flag )
// into a Version. Empty string becomes an undefined version; anything else
// becomes an exact version. This is the only way to construct a Version
// from outside the package.
func ParseVersion(s string) Version {
	if s == "" {
		return versionUndefinedValue()
	}
	return versionExactValue(s)
}

// IsDefined reports whether a concrete version was specified.
func (v Version) IsDefined() bool {
	return v.kind == versionExact
}

// String returns the raw version string. Panics if called on an undefined
// version — callers must check IsDefined() first.
func (v Version) String() string {
	if v.kind == versionUndefined {
		panic("core.Version.String: called on an undefined version")
	}
	return v.value
}












type archKind int


const (
	archUndefined archKind = iota
	archExact
)


type Arch struct {
	kind  archKind
	value string
}

func archUndefinedValue() Arch {
	return Arch{kind: archUndefined}
}

func archExactValue(v string) Arch {
	return Arch{kind: archExact, value: v}
}


func ParseArch(s string) Arch {
	if s == "" {
		return archUndefinedValue()
	}

    switch s {

    case "x86","arm64","x64":
        return archExactValue(s)

    // should not reach this point arch is already check buy cli 
    default:
        panic("unkown arch passed")
    }
	
}

func (a Arch) IsDefined() bool {
	return a.kind == archExact
}

func (a Arch) String() string {
	if a.kind == archUndefined {
		panic("core.Arch.String: called on an undefined version")
	}
	return a.value
}























/////////////////



var ErrPkgNotFound = errors.New("not found")
var ErrFetch = errors.New("fetch failed")