package install


import (
	"github.com/kasperjack/pact/core"
	"fmt"
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



func NewManaged(args core.InstallArgs, localState core.LocalState, repo core.Repo, lf core.LockFile) *installer {

	m := manager{
		localState: localState,
		repo: repo,
		lockFile: lf,
	}



	return &installer{args: args,manager: m}



}


func (i *installer) Run()error{



	i.checkNotInstalled()

	

	i.checkExists()




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


func (i *installer) checkExists() error {

	return nil
}


func (i *installer) loadBundle() error {

	return nil
}