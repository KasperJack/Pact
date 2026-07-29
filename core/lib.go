package core

import (
	//"github.com/kasperjack/pact/core/model"

	"errors"
    "runtime"
    "fmt"
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







type LocalState interface {
	// this is an fs pov
	//CreatePackage(string) error
	CreateInstallDir(PackageIdentifier string, version string) (string,error)
	PackageExists(string) (bool, error)

}
type Repo interface {

    // checks whether a package exists in the repository.
	// Returns true if the package is present, false otherwise.
	// Errors on I/O/parse failure
	PackageExists(PackageIdentifier string) (bool,error) 
                                                                
                                                              

    // loads version-independent package metadata.
    // Errors if the package doesn't exist or on I/O/parse failure
    LoadPackageInfo(packageIdentifier string) (PackageInfo,error)

     

    // resolves an upstream version (e.g. "2.7.0")
    // to its latest release for the given arch (2.7.0 ---> 2.7.0-5 --> retrun Release )
    // Errors if the package doesn't exist, the upstream version is unknown,
    // or on I/O/parse failure
    LoadReleaseByUpstreamVersion(arch Arch, packageIdentifier, upstreamVersion string) (Release, error)


    // loads an exact release by its full version
    // (e.g. "2.7.0-1"), bypassing the "latest revision" lookup.
    // Errors if the package doesn't exist, the full version is unknown,
    // or on I/O/parse failure.
    LoadReleaseByFullVersion(arch Arch, packageIdentifier, fullVersion string) (Release, error)



	//LoadPackage(PackageIdentifier string, version string, arch Arch) (*Package,ReleaseSource ,error)
	//GetVersions(PackageIdentifier string) ([]string,error)
    //GetLatest(PackageIdentifier string) (string, error)
    //GetVersionInfo(identifier, version string) (VersionInfo, error)
    //GetLatestVersionForArch(identifier string, arch Arch) (VersionInfo ,bool ,error)

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


////////////////////////



type PackageInfo struct {
    Package     string `hcl:"package"`
    Name        string `hcl:"name"`   
    Description string `hcl:"description"`
    Homepage    string `hcl:"homepage"`
    License     string `hcl:"license"`
}



type Release struct {
    Package         string `hcl:"package"`
    FullVersion     string `hcl:"version"` 
    UpstreamVersion string `hcl:"upstream_version"` 
    Revision        int     `hcl:"revision"` 
    URL             string  `hcl:"url"` 
    SHA256          string `hcl:"sha256"` 
    SizeMB          int  `hcl:"size_mb"` 
    Architecture    string `hcl:"architecture"` 
}




type ReleaseIndex struct {
    LatestVersion   string
    VersionMappings map[string]string // upstream -> full version
    Yanked          map[string]string // full version -> reason
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









type Arch int
const (
    ArchUndefined Arch = iota // explicitl: no arch was specified

    ArchX86
    ArchX64
    ArchArm64
)

func (a Arch) String() string {
    switch a {
    case ArchX86:
        return "x86"
    case ArchX64:
        return "x64"
    case ArchArm64:
        return "arm64"
    default:
        //panic("ops")
        return "Undefined"
    }
}



func ParseArch(s string) (Arch, error) {
	switch s {
	case "x86":
		return ArchX86, nil
	case "x64":
		return ArchX64, nil
	case "arm64":
		return ArchArm64, nil
	default:
		return ArchUndefined, fmt.Errorf("unknown architecture %q", s)
	}
}



func HostArch() Arch {
    switch runtime.GOARCH {
    case "386":
        return ArchX86
    case "amd64":
        return ArchX64
    default:
        return ArchArm64

    }
}












/////////////////



var ErrPkgNotFound = errors.New("not found")
var ErrFetch = errors.New("fetch failed")