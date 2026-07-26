package install

import (
	"fmt"
	//"go/version"
	"slices"
	"github.com/kasperjack/pact/core"
)


type installKind int

const (
	installKindUnknown installKind = iota
	installKindFirst              // no versions of this package installed yet
	installKindAlongside          // installing a new version next to existing ones
)


type manager struct {
    localState core.LocalState
    repo       core.Repo
    lockFile   core.LockFile
}



type installer struct {
    args          core.InstallArgs
	
	manager manager

	kind         installKind
	bundle        *core.PackageBundle
    resolvedRel   core.ReleaseSource
    resolvedArch  string
    stagingDir    string
    installDir    string
    currentDir    string
}



func NewManaged(args core.InstallArgs, localState core.LocalState, repo core.Repo, lf core.LockFile) (*installer,error) {
		//resolve version 
		//resolve arch
		// resolve install type 










	m := manager{
		localState: localState,
		repo: repo,
		lockFile: lf,
	}



	return &installer{args: args,manager: m},nil

}


func (i *installer) Run()error{
	// pcakge already exists at this point 


	err := i.checkNotInstalled()
	if err != nil {
		return nil
	}

	

	i.checkVersionExists()




	return nil
}






func (i *installer) checkNotInstalled() error {


	// pacakgeid /version 
	installed  := i.manager.lockFile.GetInstalled(i.args.PackageIdentifier)


	if i.args.Version.IsDefined() {
		if _, exists := installed[i.args.Version.String()]; exists {
			return fmt.Errorf("package %q version %q already installed",
				i.args.PackageIdentifier, i.args.Version.String())
		}

		if len(installed) > 0 {
			i.kind = installKindAlongside
			// installing a version along side another version
		} else {
			i.kind = installKindFirst
		}



		return nil
	}
	

	if len(installed) > 0 {
		return fmt.Errorf(
			"package %q already installed with other version(s); specify a version to install alongside them",
			i.args.PackageIdentifier,
		)
	}



	i.kind = installKindFirst
	return nil // first install — version will be resolved to latest downstream

}


func (i *installer) checkVersionExists() error {

	if i.args.Version.IsDefined() {


	}



	return nil
}








func (i *installer) loadBundle() error {

	return nil
}




func resolveVersion(version core.Version, packageIdentifier string, repo core.Repo)(string,error){

			if version.IsDefined() {
			//hanndle errors //fetch //should be sfve to asume not no pkgnotfound error
			versions,_ := repo.GetVersions(packageIdentifier)

			if !slices.Contains(versions,version.String()) {
				return "",fmt.Errorf("pkg %s does not have version %s",packageIdentifier,version.String())
			}
			return version.String(),nil
		}

		return "",nil
}