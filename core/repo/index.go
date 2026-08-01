package repo

import(
	"path/filepath"
	"os"
	"github.com/kasperjack/pact/core"
    //"github.com/kasperjack/pact/core/parce"
	//"fmt"
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