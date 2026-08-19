package repo

import(
	//"path/filepath"
	//"os"
	//"github.com/kasperjack/pact/core"
    //"github.com/kasperjack/pact/core/parce"
	//"fmt"
    //"slices"
)


/*
type Package struct {
    repo *repo
    arch core.Arch
    id   string
	index core.ReleaseIndex
}






func (p *Package) Info() (core.PackageInfo, error) {

	return p.repo.LoadPackageInfo(p.id)
}



func (p *Package) LatestReleaseVersion() string {
    return p.index.LatestVersion
}



func (p *Package) LatestUpstreamVersion() string {
    return p.index.UpstreamOf[p.index.LatestVersion]
}





func (p *Package) ReleaseVersions() []string {
    versions := make([]string, 0, len(p.index.UpstreamOf))

    for version := range p.index.UpstreamOf {
        versions = append(versions, version)
    }

    return versions
}

func (p *Package) UpstreamVersions() []string {
    versions := make([]string, 0, len(p.index.VersionMappings))

    for version := range p.index.VersionMappings {
        versions = append(versions, version)
    }

    return versions
}

func (p *Package) LatestReleaseForUpstream(upstreamVersion string) string {

    return p.index.VersionMappings[upstreamVersion]
}




func (p *Package) Release(releaseVersion string) (core.Release, error){

	return p.repo.LoadRelease(p.arch,p.id,releaseVersion)
}
*/