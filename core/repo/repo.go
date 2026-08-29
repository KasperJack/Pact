package repo


import(
	"path/filepath"
	"os"
	"github.com/kasperjack/pact/core"
    "github.com/kasperjack/pact/core/parce"
	"fmt"
    //"slices"
	"strconv"
)

type repo struct {
	repoRoot string
}

// establish erly that pkg exits and arch is valid and version exists 

func NewLocalRepo(repoRoot string) (core.Repo,error) {

    _, err := os.Stat(repoRoot)

    if err != nil {

        return nil,fmt.Errorf("can't find the test bucket") //RF:E
    }


	return &repo{repoRoot: repoRoot},nil
}













func (r *repo) PackageExists(packageIdentifier string) (bool, error) {
	dirPath := filepath.Join(r.repoRoot, "packages",packageIdentifier)

	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}
















func (r *repo) LoadPackageInfo(packageIdentifier string) (*core.PackageInfo,error) {

    pkgFilePath := filepath.Join(r.repoRoot,"packages",packageIdentifier,"package.hcl")



    pkgData, err := os.ReadFile(pkgFilePath)
    if err != nil {
        return nil,fmt.Errorf("%w: loading %s fn:LoadPackageInfo: %v", core.ErrFetch, packageIdentifier, err)
    }

	p, diags := parce.PackageInfo(pkgData)

	if diags.HasErrors(){
		return  nil,diags
	}



    return p,nil

}






func (r *repo) LoadPackageIndex(packageIdentifier string) (core.PackageIndex, error) {


	indexFilePath := filepath.Join(r.repoRoot, "packages", packageIdentifier, "index.hcl")

	indexData, err := os.ReadFile(indexFilePath)
	if err != nil {
		return core.PackageIndex{}, fmt.Errorf("%w: loading index %s fn:LoadPackageIndex: %v", core.ErrFetch, packageIdentifier, err)
	}

	index, err := parce.PackageIndex(indexData)
	if err != nil {
		return core.PackageIndex{}, fmt.Errorf("%w: parsing index %s fn:LoadPackageIndex: %v", core.ErrFetch, packageIdentifier, err)
	}



	return index,nil

}







func (r *repo) LoadArchRelease(packageIdentifier, Version string, arch core.Arch, revision int) (*core.ArchRelease, error) {


	manifetFilePath := filepath.Join(r.repoRoot, "packages",packageIdentifier, Version, arch.String() ,strconv.Itoa(revision), "manifest.hcl")

	data, err := os.ReadFile(manifetFilePath)
	if err != nil {
		return nil, fmt.Errorf("%w: loading manifest %s@%s for arch %s fn:LoadArchRelease: %v", core.ErrFetch, packageIdentifier, Version, arch.String(), err)
	}


	m, diags :=  parce.Manifest(data)

	if diags.HasErrors() {
		return nil, diags
	}




	releaseFilePath := filepath.Join(r.repoRoot, "packages",packageIdentifier, Version,arch.String() ,strconv.Itoa(revision), "release.hcl")


	data, err = os.ReadFile(releaseFilePath)

	if err != nil {
		return nil, fmt.Errorf("%w: loading release %s@%s for arch %s fn:LoadArchRelease: %v", core.ErrFetch, packageIdentifier, Version, arch.String(), err)
	}

	re, err := parce.Release(data)

	if err != nil {

		return nil,fmt.Errorf("%w: parsing release %s fn:LoadArchRelease: %v", core.ErrFetch, packageIdentifier, err)
	}

	return &core.ArchRelease{
		Manifest: m,
		Release: re,
		Interface: core.Interface{},

	},nil
}








func (r *repo) LoadArchStatus(packageIdentifier, Version string, arch core.Arch) (core.ArchStatus, error) {

	statusFilePath := filepath.Join(r.repoRoot, "packages", packageIdentifier, Version, arch.String(), "status.hcl")


	data,err := os.ReadFile(statusFilePath)
	if err != nil {
		return core.ArchStatus{}, fmt.Errorf("%w: loading status %s fn:LoadArchStatus: %v", core.ErrFetch, packageIdentifier, err)
	}


	return parce.ArchStatus(data)


}



















































/*



func (r *repo) Package(arch core.Arch,packageIdentifier string) (core.Package, error) {

    ok, err := r.PackageExists(packageIdentifier)
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

    index, err := r.LoadPackageIndex(packageIdentifier)
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