package main

import(
	"path"
	"os"
	"github.com/kasperjack/pact/core"
    "github.com/kasperjack/pact/core/parce"
	"fmt"
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






func (r *repo) PackageExists(packageName string) (bool,error) {

	pkgPath := path.Join(r.repoRoot,packageName)

	info, err := os.Stat(pkgPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil 
		}
    	return false, err  // permission or something else, real error
	}

	if  !info.IsDir() {
		return false,nil
	}

	return true,nil
}



func (r *repo) GetVersions(packageName string) ([]string, error) {


    pkgPath := path.Join(r.repoRoot, packageName)

    entries, err := os.ReadDir(pkgPath)
    if err != nil {
        return nil, err
    }

    var versions []string
    for _, entry := range entries {
        if entry.IsDir() {
            versions = append(versions, entry.Name())
        }
    }

    return versions, nil
}




func (r *repo) LoadPackage(packageName string, version string) (core.PackageBundle,error) {
    
    pkgFilePath := path.Join(r.repoRoot, packageName, fmt.Sprintf("%s.hcl", packageName))
    pkgData, err := os.ReadFile(pkgFilePath)
    if err != nil {
        return core.PackageBundle{}, err
    }
    pkg, err := parce.Pacakge(pkgData)
    if err != nil {
        return core.PackageBundle{}, err
    }

    // read and parse release manifest
    releaseFilePath := path.Join(r.repoRoot, packageName, version, "release.hcl") //RF:E (can't find release file for version)
    releaseData, err := os.ReadFile(releaseFilePath)
    if err != nil {
        return core.PackageBundle{}, err
    }
    release, err := parce.Release(releaseData)
    if err != nil {
        return core.PackageBundle{}, err
    }

    // lua script just stays as raw bytes, core will run it
    scriptFilePath := path.Join(r.repoRoot, packageName, version, "script.star") //RF:E (can't find install script )
    script, err := os.ReadFile(scriptFilePath)
    if err != nil {
        return core.PackageBundle{}, err
    }

	f := core.PackageBundle{
		Release: release,
		Package: pkg,
		Script: script,
	}



    return f,nil
}