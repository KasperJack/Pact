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


    return core.Release{},nil
}







func (r *repo) LoadReleaseByFullVersion(arch core.Arch, packageIdentifier, fullVersion string) (core.Release, error) {



    return core.Release{},nil
}










