package repo


import(
	"path/filepath"
	"os"
	"github.com/kasperjack/pact/core"
    "github.com/kasperjack/pact/core/parce"
	"fmt"
    //"slices"
)

type repo struct {
	repoRoot string
}



func NewLocalRepo(repoRoot string) (core.Repo,error) {

    _, err := os.Stat(repoRoot)

    if err != nil {

        return nil,fmt.Errorf("can't find the test bucket") //RF:E
    }


	return &repo{repoRoot: repoRoot},nil
}



func (r *repo) Package(arch core.Arch,packageIdentifier string) (core.Package, error) {

    ok, err := r.HasPackage(arch, packageIdentifier)
    if err != nil {
        return nil, err
    }

    if !ok {
        return nil, fmt.Errorf(
            "%w: %s for arch %s",
            core.ErrPkgNotFound,
            packageIdentifier,
            arch.String(),
        )
    }

    index, err := r.LoadIndex(arch, packageIdentifier)
    if err != nil {
        return nil, err
    }

    return &Package{
        repo:  r,
        arch:  arch,
        id:    packageIdentifier,
        index: index,
    }, nil
}








func (r *repo) HasPackage(arch core.Arch, packageIdentifier string) (bool, error) {
	dirPath := filepath.Join(r.repoRoot, "releases", arch.String(), packageIdentifier)

	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}





func (r *repo) LoadPackageInfo(packageIdentifier string) (core.PackageInfo,error) {

    pkgFilePath := filepath.Join(r.repoRoot,"packages",packageIdentifier,"package.hcl")



    pkgData, err := os.ReadFile(pkgFilePath)
    if err != nil {
        return core.PackageInfo{},fmt.Errorf("%w: loading %s fn:LoadPackageInfo: %v", core.ErrFetch, packageIdentifier, err)
    }

    return parce.PackageInfo(pkgData)

}






func (r *repo) LoadIndex(arch core.Arch, packageIdentifier string) (core.ReleaseIndex, error) {


	indexFilePath := filepath.Join(r.repoRoot, "releases", arch.String(), packageIdentifier, "index.hcl")

	indexData, err := os.ReadFile(indexFilePath)
	if err != nil {
		return core.ReleaseIndex{}, fmt.Errorf("%w: loading index %s fn:LoadReleaseByUpstreamVersion: %v", core.ErrFetch, packageIdentifier, err)
	}

	index, err := parce.ReleaseIndex(indexData)
	if err != nil {
		return core.ReleaseIndex{}, fmt.Errorf("%w: parsing index %s fn:LoadReleaseByUpstreamVersion: %v", core.ErrFetch, packageIdentifier, err)
	}



	return index,nil

}







func (r *repo) LoadRelease(arch core.Arch, packageIdentifier, fullVersion string) (core.Release, error) {


	pkgFile := filepath.Join(r.repoRoot, "releases", arch.String(), packageIdentifier, fullVersion, "package.hcl")

	fullVersionData, err := os.ReadFile(pkgFile)
	if err != nil {
		return core.Release{}, fmt.Errorf("%w: loading package %s@%s fn:LoadReleaseByFullVersion: %v", core.ErrFetch, packageIdentifier, fullVersion, err)
	}

	return parce.Release(fullVersionData)


}






















































/*




func (r *repo) LoadReleaseByUpstreamVersion(arch core.Arch, packageIdentifier, upstreamVersion string) (core.Release, error) {



	indexFilePath := filepath.Join(r.repoRoot, "releases", arch.String(), packageIdentifier, "index.hcl")

	indexData, err := os.ReadFile(indexFilePath)
	if err != nil {
		return core.Release{}, fmt.Errorf("%w: loading index %s fn:LoadReleaseByUpstreamVersion: %v", core.ErrFetch, packageIdentifier, err)
	}

	index, err := parce.ReleaseIndex(indexData)
	if err != nil {
		return core.Release{}, fmt.Errorf("%w: parsing index %s fn:LoadReleaseByUpstreamVersion: %v", core.ErrFetch, packageIdentifier, err)
	}

	fullVersion, ok := index.VersionMappings[upstreamVersion]
	if !ok {
		return core.Release{}, fmt.Errorf("%w: unknown upstream version %q for %s/%s fn:LoadReleaseByUpstreamVersion", core.ErrVersionNotFound, upstreamVersion, arch.String(),packageIdentifier)
	}

	return r.LoadReleaseByFullVersion(arch, packageIdentifier, fullVersion)
}




func (r *repo) LoadReleaseByFullVersion(arch core.Arch, packageIdentifier, fullVersion string) (core.Release, error) {


    ok, err := r.PackageExists(arch,packageIdentifier)
    if err != nil {
		return core.Release{}, err
	}

    if !ok {
		return core.Release{}, fmt.Errorf("%w: %s for arch %s fn:LoadReleaseByFullVersion", core.ErrPackageNotFoundForArch, packageIdentifier, arch.String()) //pkgnotfoundfor arch
	}




	pkgFile := filepath.Join(r.repoRoot, "releases", arch.String(), packageIdentifier, fullVersion, "package.hcl")

	fullVersionData, err := os.ReadFile(pkgFile)
	if err != nil {
		return core.Release{}, fmt.Errorf("%w: loading package %s@%s fn:LoadReleaseByFullVersion: %v", core.ErrFetch, packageIdentifier, fullVersion, err)
	}

	return parce.Release(fullVersionData)
}



// get upstream versions 
// get full versions 


func (r *repo) GetLatestFullVersion(arch core.Arch, packageIdentifier string) (string, error) {




    ok, err := r.PackageExists(arch,packageIdentifier)
    if err != nil {
		return "", err
	}

    if !ok {
        return "", fmt.Errorf("%w: %s for arch %s fn:LoadLatestFullVersion", core.ErrPackageNotFoundForArch, packageIdentifier, arch.String())
    }


	indexFilePath := filepath.Join(r.repoRoot, "releases", arch.String(), packageIdentifier, "index.hcl")

	indexData, err := os.ReadFile(indexFilePath)
	if err != nil {
		return "", fmt.Errorf("%w: loading index %s fn:LoadLatestVersion: %v", core.ErrFetch, packageIdentifier, err)
	}

	index, err := parce.ReleaseIndex(indexData)
	if err != nil {
		return "", fmt.Errorf("%w: parsing index %s fn:LoadLatestVersion: %v", core.ErrFetch, packageIdentifier, err)
	}

	if index.LatestVersion == "" {
		return "", fmt.Errorf("%w: no latest_version set for %s fn:LoadLatestVersion", core.ErrVersionNotFound, packageIdentifier)
	}

	return index.LatestVersion, nil
}





*/