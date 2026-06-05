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



func NewLocalRepo(repoRoot string) core.Repo {

	return &repo{repoRoot: repoRoot}
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




func (r *repo) LoadPackage(packageName string, version string) (core.PackageFiles,error) {

    pkgFilePath := path.Join(r.repoRoot, packageName, fmt.Sprintf("%s.hcl", packageName))
    pkgData, err := os.ReadFile(pkgFilePath)
    if err != nil {
        return core.PackageFiles{}, err
    }
    pkg, err := parce.Pacakge(pkgData)
    if err != nil {
        return core.PackageFiles{}, err
    }

    // read and parse release manifest
    releaseFilePath := path.Join(r.repoRoot, packageName, version, "release.hcl")
    releaseData, err := os.ReadFile(releaseFilePath)
    if err != nil {
        return core.PackageFiles{}, err
    }
    release, err := parce.Release(releaseData)
    if err != nil {
        return core.PackageFiles{}, err
    }

    // lua script just stays as raw bytes, core will run it
    scriptFilePath := path.Join(r.repoRoot, packageName, version, "script.lua")
    script, err := os.ReadFile(scriptFilePath)
    if err != nil {
        return core.PackageFiles{}, err
    }

	f := core.PackageFiles{
		Release: release,
		Package: pkg,
		LuaScript: script,
	}



    return f,nil
}