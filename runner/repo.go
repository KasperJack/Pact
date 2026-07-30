package main

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
    } // check if is a a valid pact repo . ? 


	return &repo{repoRoot: repoRoot},nil
}





//index //fetch error
func (r *repo) PackageExists(packageIdentifier string) (bool,error) { //add type of package 

	pkgFilePath := filepath.Join(r.repoRoot,"packages",packageIdentifier,"package.hcl")


	_, err := os.Stat(pkgFilePath)

	if err != nil {
		if os.IsNotExist(err) {
			return false,nil 
		}
    	return false,fmt.Errorf("%w: loading :%s : %v", core.ErrFetch, packageIdentifier, err)  // permission or something else
	}

	return true,nil
}




//index ? 
func (r *repo) LoadPackageInfo(packageIdentifier string) (core.PackageInfo,error) {

    pkgFilePath := filepath.Join(r.repoRoot,"packages",packageIdentifier,"package.hcl")

    ok, err := r.PackageExists(packageIdentifier)

    if err !=nil {
        return core.PackageInfo{},err
    }
    if !ok {return core.PackageInfo{},fmt.Errorf("pkg %s not found fn:LoadPackageInfo",packageIdentifier)}


    pkgData, err := os.ReadFile(pkgFilePath)
    if err != nil {
        return core.PackageInfo{},fmt.Errorf("%w: loading %s fn:LoadPackageInfo: %v", core.ErrFetch, packageIdentifier, err)
    }

    return parce.PackageInfo(pkgData)

}



func (r *repo) LoadReleaseByUpstreamVersion(arch core.Arch, packageIdentifier, upstreamVersion string) (core.Release, error) {




    ok, err := r.PackageExistsForArch(arch,packageIdentifier)
    if err != nil {
		return core.Release{}, err
	}

    if !ok {
		return core.Release{}, fmt.Errorf("%w: %s for arch %s fn:LoadReleaseByUpstreamVersion", core.ErrPackageNotFoundForArch, packageIdentifier, arch.String())
	}



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
		return core.Release{}, fmt.Errorf("%w: unknown upstream version %q for %s fn:LoadReleaseByUpstreamVersion", core.ErrVersionNotFound, upstreamVersion, packageIdentifier)
	}

	return r.LoadReleaseByFullVersion(arch, packageIdentifier, fullVersion)
}




func (r *repo) LoadReleaseByFullVersion(arch core.Arch, packageIdentifier, fullVersion string) (core.Release, error) {


    ok, err := r.PackageExistsForArch(arch,packageIdentifier)
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






func (r *repo) LoadLatestFullVersion(arch core.Arch, packageIdentifier string) (string, error) {




    ok, err := r.PackageExistsForArch(arch,packageIdentifier)
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









func (r *repo) PackageExistsForArch(arch core.Arch, packageIdentifier string) (bool, error) {
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