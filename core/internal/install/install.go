package install


import (
	"github.com/kasperjack/pact/core"
	"fmt"
)


type Installer interface{
	Run() error
}



type manager struct {
    localState core.LocalState
    repo       core.Repo
    lockFile   core.LockFile
}



type installer struct {
    args          core.InstallArgs
	
	manager manager


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


	if i.args.Version == "" {

		if len(installed) > 0 {
			// error: pkg already installed must provied a spesefic to install 
		}
		
		
		// first install will install latest 
	}else{  // a version was passed 

		if _, exists := installed[i.args.Version]; exists {
    		return fmt.Errorf("package %q version %q already installed", i.args.PackageIdentifier, i.args.Version)
		}
		// installing pkg@version
	}
	








	return nil
}


func (i *installer) checkExists() error {

	return nil
}


func (i *installer) loadBundle() error {

	return nil
}