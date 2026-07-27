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
func (r *repo) PackageExists(PackageIdentifier string) (bool,string,error) { //add type of package 

	pkgFilePath := filepath.Join(r.repoRoot,PackageIdentifier,fmt.Sprintf("%s.hcl",PackageIdentifier))


    //fmt.Println(pkgFilePath)

	_, err := os.Stat(pkgFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, "",nil 
		}
    	return false,"", fmt.Errorf("%w: fetching %s: %v", core.ErrFetch, PackageIdentifier, err)  // permission or something else
	}

    pkgData, err := os.ReadFile(pkgFilePath)
    if err != nil {
        return false, "",fmt.Errorf("%w: fetching %s: %v", core.ErrFetch, PackageIdentifier, err)
    }


    pType,err := parce.GetType(pkgData)

    if err != nil{
        return false,"",fmt.Errorf("%w: fetching %s: %v", core.ErrFetch, PackageIdentifier, err)
    }


	return true,pType,nil
}

            
//index  //fetch error //not found error 
func (r *repo) GetVersions(PackageIdentifier string) ([]string, error) { //error ftching 


    indexFilePath := filepath.Join(r.repoRoot,PackageIdentifier,"index.hcl")


    _, err := os.Stat(indexFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil,core.ErrPkgNotFound
		}
    	return nil, fmt.Errorf("%w: fetching versions %s: %v", core.ErrFetch, PackageIdentifier, err)  // permission or something else
	}



    pkgData, err := os.ReadFile(indexFilePath)
    if err != nil {
        return nil,fmt.Errorf("%w: fetching versions%s: %v", core.ErrFetch, PackageIdentifier, err)
    }


    versions,err := parce.GetVersions(pkgData)
    if err != nil{
        return nil,fmt.Errorf("%w: fetching versions %s: %v", core.ErrFetch, PackageIdentifier, err)
    }



    return versions, nil
}









func (r *repo) LoadPackage(PackageIdentifier string, version string) (core.PackageBundle,error) {
    
    pkgFilePath := filepath.Join(r.repoRoot, PackageIdentifier, fmt.Sprintf("%s.hcl", PackageIdentifier))
    pkgData, err := os.ReadFile(pkgFilePath)
    if err != nil {
        return core.PackageBundle{}, err
    }
    pkg, err := parce.Pacakge(pkgData)
    if err != nil {
        return core.PackageBundle{}, err
    }

    // read and parse release manifest
    releaseFilePath := filepath.Join(r.repoRoot, PackageIdentifier, version, "release.hcl") //RF:E (can't find release file for version)
    releaseData, err := os.ReadFile(releaseFilePath)
    if err != nil {
        return core.PackageBundle{}, err
    }
    release, err := parce.Release(releaseData)
    if err != nil {
        return core.PackageBundle{}, err
    }

    // lua script just stays as raw bytes, core will run it
    scriptFilePath := filepath.Join(r.repoRoot, PackageIdentifier, version, "script.star") //RF:E (can't find install script )
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


// index
func (r *repo) GetVersionInfo(PackageIdentifier, version string) (core.VersionInfo, error) {

    releaseFilePath := filepath.Join(r.repoRoot,PackageIdentifier,version,"release.hcl")


    _, err := os.Stat(releaseFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return core.VersionInfo{},core.ErrPkgNotFound
		}
    	return core.VersionInfo{}, fmt.Errorf("%w: fetching version info %s: %v", core.ErrFetch, PackageIdentifier, err)  
	}


    pkgData, err := os.ReadFile(releaseFilePath)
    if err != nil {
        return core.VersionInfo{},fmt.Errorf("%w: fetching version info %s: %v", core.ErrFetch, PackageIdentifier, err)
    }


    arches,err := parce.GetSourceBlockNames(pkgData)

    if err != nil{
    return core.VersionInfo{},fmt.Errorf("%w: fetching version info %s: %v", core.ErrFetch, PackageIdentifier, err)
    }

    vi := core.VersionInfo{ArchFallbackSafe: true,Version: version} // safe is the defult adding  a key  to signal unsafe 
   
	for _, a := range arches {     

		arch, err := core.ParseArch(a)
		if err != nil {
			return core.VersionInfo{}, fmt.Errorf("parsing arch: %w", err)
		}
		vi.Archs = append(vi.Archs, arch) // mutating a field ON vi
	}

    return vi,nil

}



func (r *repo) GetLatestVersionForArch(PackageIdentifier string, arch core.Arch) (core.VersionInfo ,bool ,error){

    indexFilePath := filepath.Join(r.repoRoot,PackageIdentifier,"index.hcl")


    _, err := os.Stat(indexFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return core.VersionInfo{},false,core.ErrPkgNotFound
		}
    	return core.VersionInfo{},false, fmt.Errorf("%w: fetching version for arch %s: %v", core.ErrFetch, PackageIdentifier, err)  // permission or something else
	}



    indexData, err := os.ReadFile(indexFilePath)
    if err != nil {
        return core.VersionInfo{},false,fmt.Errorf("%w: fetching version for arch %s: %v", core.ErrFetch, PackageIdentifier, err)
    }



     v,ok,err := parce.GetLatestVerArch(indexData,arch.String())

    if err != nil {

        return core.VersionInfo{},false,err
    }

    if ok {

        return core.VersionInfo{Version: v,ArchFallbackSafe: true},true,nil

    }


    return core.VersionInfo{},false,nil
}





















/*
//index //fetch error  //not found error 
func (r *repo) GetLatest(PackageIdentifier string) (string, error) { //error ftching 


    pkgFilePath := filepath.Join(r.repoRoot,PackageIdentifier,"index.hcl")


    _, err := os.Stat(pkgFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "",core.ErrPkgNotFound
		}
    	return "", fmt.Errorf("%w: fetching latest %s: %v", core.ErrFetch, PackageIdentifier, err)  
	}


    pkgData, err := os.ReadFile(pkgFilePath)
    if err != nil {
        return "",fmt.Errorf("%w: fetching latest %s: %v", core.ErrFetch, PackageIdentifier, err)
    }


    latest,err := parce.GetLatest(pkgData)

        if err != nil{
        return "",fmt.Errorf("%w: fetching latest %s: %v", core.ErrFetch, PackageIdentifier, err)
    }

    return latest, nil
}

*/