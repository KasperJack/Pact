package main

import(
	"path/filepath"
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


    pkgFilePath := filepath.Join(r.repoRoot,PackageIdentifier,fmt.Sprintf("%s.hcl",PackageIdentifier))


    _, err := os.Stat(pkgFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil,core.ErrPkgNotFound
		}
    	return nil, fmt.Errorf("%w: fetching %s: %v", core.ErrFetch, PackageIdentifier, err)  // permission or something else
	}



    pkgData, err := os.ReadFile(pkgFilePath)
    if err != nil {
        return nil,fmt.Errorf("%w: fetching %s: %v", core.ErrFetch, PackageIdentifier, err)
    }


    versions,err := parce.GetVersions(pkgData)
    if err != nil{
        return nil,fmt.Errorf("%w: fetching %s: %v", core.ErrFetch, PackageIdentifier, err)
    }



    return versions, nil
}


//index //fetch error  //not found error 
func (r *repo) GetLatest(PackageIdentifier string) (string, error) { //error ftching 


    pkgFilePath := filepath.Join(r.repoRoot,PackageIdentifier,fmt.Sprintf("%s.hcl",PackageIdentifier))


    _, err := os.Stat(pkgFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "",core.ErrPkgNotFound
		}
    	return "", fmt.Errorf("%w: fetching %s: %v", core.ErrFetch, PackageIdentifier, err)  // permission or something else
	}


    pkgData, err := os.ReadFile(pkgFilePath)
    if err != nil {
        return "",fmt.Errorf("%w: fetching %s: %v", core.ErrFetch, PackageIdentifier, err)
    }


    latest,err := parce.GetLatest(pkgData)

        if err != nil{
        return "",fmt.Errorf("%w: fetching %s: %v", core.ErrFetch, PackageIdentifier, err)
    }

    return latest, nil
}




// loadLatest using tags from indexes 

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