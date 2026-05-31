package main

import(
	"path"
	"os"
	"github.com/kasperjack/pact/core"
	"fmt"
)

type Repo struct {
	RepoRoot string
}



func (r *Repo) PackageExists(packageName string) (bool,error) {

	pkgPath := path.Join(r.RepoRoot,packageName)

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



func (r *Repo) GetVersions(packageName string) ([]string, error) {


    pkgPath := path.Join(r.RepoRoot, packageName)

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




func (r *Repo) LoadPackage(packageName string, version string) (core.PackageFiles,error) {

	pkgPath := path.Join(r.RepoRoot,packageName)



	pkgFilePath :=  path.Join(pkgPath,fmt.Sprintf("%s.hcl",packageName))
	pkgFile, err := os.Open(pkgFilePath)
	if err != nil {
		return core.PackageFiles{},err
	}
	defer pkgFile.Close()



	releaseFilePath := path.Join(pkgPath,version,"release.hcl")
	releaseFile, err := os.Open(releaseFilePath)
	if err != nil {
		return core.PackageFiles{},err
	}
	defer releaseFile.Close()

	sciptFilePath := path.Join(pkgPath,version,"script.lua")
	scriptFile, err := os.Open(sciptFilePath)
	if err != nil {
		return core.PackageFiles{},err
	}
	defer scriptFile.Close()

	pf := core.PackageFiles{
		Package: pkgFile,
		Release: releaseFile,
		LuaScript: scriptFile,
		
	}

	return pf,nil
}
