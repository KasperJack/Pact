package repo

import(
	"path/filepath"
	"os"
	"github.com/kasperjack/pact/core"
    //"github.com/kasperjack/pact/core/parce"
	"fmt"
    //"slices"
)
func (r *repo) PackagHeExistsForArch(arch core.Arch, packageIdentifier string) (bool, error) {
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






//index //fetch error
func (r *repo) OPackageExists(packageIdentifier string) (bool,error) { 

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