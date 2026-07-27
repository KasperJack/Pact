package install

import (
	"fmt"
	//"go/version"
	"slices"

	"github.com/kasperjack/pact/core"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
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




func resolveArchVersion(args core.InstallArgs, repo core.Repo)(string,core.Arch,error){

		if args.Version.IsDefined() {
			//hanndle errors //fetch //should be sfve to asume not no pkgnotfound error
			versions,_ := repo.GetVersions(args.PackageIdentifier)

			if !slices.Contains(versions,args.Version.String()) {
				return "",core.ArchUndefined,fmt.Errorf("pkg %s does not have version %s",args.PackageIdentifier,args.Version.String())
			}

			// if the pkg no arch 

			archInfo,err := repo.GetVersionInfo(args.PackageIdentifier,args.Version.String())
			if err != nil {return "",core.ArchUndefined,err}
			
			switch core.HostArch() {
				
			
				case core.ArchX64:

					if args.TargetArch == core.ArchUndefined {
						if slices.Contains(archInfo.Archs,core.ArchX64){
							return args.TargetArch.String(),core.ArchX64,nil
						}

						if archInfo.ArchFallbackSafe {
							if slices.Contains(archInfo.Archs,core.ArchX64){
								return args.TargetArch.String(),core.ArchX64,nil
							}
						}

						//error

					}

					if args.TargetArch == core.ArchX86 {
						if slices.Contains(archInfo.Archs,core.ArchX64){
							return args.TargetArch.String(),core.ArchX64,nil
						}

						// error 
					}

					if slices.Contains(archInfo.Archs,core.ArchX86){
						return args.TargetArch.String(),core.ArchX86,nil
					}

					// error 



				case core.ArchX86:

					if slices.Contains(archInfo.Archs,core.ArchX86){
							return args.TargetArch.String(),core.ArchX86,nil
					}

					// error 

				default: // arm64
					if slices.Contains(archInfo.Archs,core.ArchArm64){
						return args.TargetArch.String(),core.ArchArm64,nil
					}

					// error 

			}

		}else{ // no defined verions

			if args.TargetArch == core.ArchUndefined {

				ok,v,err := repo.GetLatestVersion(args.PackageIdentifier,core.HostArch()) // get latest verion for an arch 

				if err != nil {
					return "",core.ArchUndefined,err
				}

				if ok {
					return v.verion,core.HostArch(),nil
				}

				if core.HostArch() == core.ArchArm64 {

					ok,v,_ := repo.GetLatestVersion(args.PackageIdentifier,core.ArchX86)

					if ok {
						if v.ArchFallbackSafe {
							return v.verion,core.ArchX86,nil
						}
					}

				}

				// error 
			}





			switch core.HostArch() {

				case core.ArchX64:

					if args.TargetArch == core.ArchX64 {
						ok,v,_ := repo.GetLatestVersion(args.PackageIdentifier,core.ArchX64)
					}

					if args.TargetArch == core.ArchX86 {
						
					}

				case core.ArchX86:
					ok,v,_ := repo.GetLatestVersion(args.PackageIdentifier,core.ArchX86)


				default: //arm64
					ok,v,_ := repo.GetLatestVersion(args.PackageIdentifier,core.ArchX86)
					

			}



		}











		return "",core.ArchUndefined,fmt.Errorf("unexpecteed error happend resoving the version/arch")
}







