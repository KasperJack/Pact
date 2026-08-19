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
	Desktop() string
	CacheDir() string
	InstallDir() string
	Repo() string
	LockFile() string
}









type Repo interface {

    LoadPackageInfo(packageIdentifier string) (PackageInfo,error)
    PackageExists(packageIdentifier string) (bool, error)
    LoadPackageIndex(packageIdentifier string) (PackageIndex, error)
    LoadArchStatus(packageIdentifier, Version string, arch Arch) (ArchStatus, error)
    LoadArchRelease(packageIdentifier, Version string, arch Arch, revision int) (Release, error)
}

type Package interface {

    Info() (PackageInfo, error)
    LatestReleaseVersion() string
    LatestUpstreamVersion() string 
    ReleaseVersions() []string
    UpstreamVersions() []string
    LatestReleaseForUpstream(upstreamVersion string) string
    Release(releaseVersion string) (Release, error)
}




type LockManager interface {
    
    GetInstalled(packageIdentifier string) (LockedPackage, bool)

    RecordInstall(pkg LockedPackage) error

    RecordRemove(packageIdentifier string) error

	Test() error

}
    //InstallDir(pkg string) string
	//IsInstalled(pkg string) error







////////



type LockedPackage struct {
    Identifier        string `hcl:"identifier,label"`
    Version     string `hcl:"version"`
    UpstreamVersion string `hcl:"upstream_version"`
    InstalledAt string `hcl:"installed_at"`
    InstallDir  string `hcl:"install_dir"`
    Arch        string `hcl:"arch"`
}

type LockFileC struct {
    Packages []LockedPackage `hcl:"package,block"`
}





////////////////////////




type PackageInfo struct {
    Package     string `hcl:"package"`
    Name        string `hcl:"name"`   
    Description string `hcl:"description"`
    Homepage    string `hcl:"homepage"`
    License     string `hcl:"license"`
    ArchitecturesRaw []string `hcl:"architectures"`
    Architectures   []Arch 

}



func (p *PackageInfo) Validate() error {
    p.Architectures = make([]Arch, 0, len(p.ArchitecturesRaw))

    for _, raw := range p.ArchitecturesRaw {
        arch, err := ParseArch(raw)
        if err != nil {
            return fmt.Errorf("invalid architecture %q: %w", raw, err)
        }

        p.Architectures = append(p.Architectures, arch)
    }

    return nil
}











type PackageIndex struct {
    LatestVersion string   `hcl:"latest_version"`
    Versions      []string `hcl:"versions"`
    ArchMap       ArchMap  `hcl:"arch_map,block"`
}

type ArchMap struct {
    Archs []ArchInfo `hcl:"arch,block"`
}

type ArchInfo struct {
    ArchRaw     string `hcl:"name,label"`
    Arch       Arch //RT:V
    Version  string `hcl:"version"`
    Revision int    `hcl:"revision"`
}


func (p *PackageIndex) Validate() error {

    for i := range p.ArchMap.Archs {

        a, err := ParseArch(p.ArchMap.Archs[i].ArchRaw)


        if err != nil {
            return fmt.Errorf("index: %w", err)
        }

        p.ArchMap.Archs[i].Arch = a

    }


    return nil
}





func (a ArchMap) AsMap() map[Arch]ArchInfo {
    m := make(map[Arch]ArchInfo, len(a.Archs))
    for _, arch := range a.Archs {

        m[arch.Arch] = arch
    }
    return m
}





type ArchStatus struct {
    Status string `hcl:"status"`
    ReleasedAt string `hcl:"released_at,optional"`
    CurrentRevision int `hcl:"current_revision,optional"`
}











type Release struct {
    Package         string `hcl:"package"`
    UpstreamVersion string `hcl:"upstream_version"` 
    Revision        int     `hcl:"revision"` 
    URL             string  `hcl:"url"` 
    SHA256          string `hcl:"sha256"` 
    SizeMB          int  `hcl:"size_mb"` 

    ArchitectureRaw string `hcl:"architecture"`
    
    Architecture Arch //RT:V
}


func (r *Release) Validate() error {
    a, err := ParseArch(r.ArchitectureRaw)

    if err != nil {
        return fmt.Errorf("release %q: %w", r.Package, err)
    }
    r.Architecture = a 
    return nil
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



func (a Arch) CompatibleArchs() []Arch {
    switch a {
    case ArchArm64:

        return []Arch{ArchArm64}

    case ArchX64:
        return []Arch{ArchX64, ArchX86,}

    case ArchX86:
        return []Arch{ArchX86}

    default:
        return nil
    }
}

func (a Arch) Priority(target Arch) (int, bool) {
    for i, arch := range a.CompatibleArchs() {
        if arch == target {
            return i, true
        }
    }
    return 0, false
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

var (

    ErrPkgNotFound = errors.New("pkg not found")
    ErrFetch = errors.New("fetch failed")
    ErrVersionNotFound = errors.New("version not found")
    ErrPackageNotFound  = errors.New("package not found")
    ErrPackageNotFoundForArch = errors.New("package not found for arch")
)

