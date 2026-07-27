package install

import (
	"fmt"
	//"go/version"
	"slices"

	"github.com/kasperjack/pact/core"
	//"golang.org/x/tools/go/analysis/passes/nilfunc"
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




		ver,arch, err := resolveArchVersion(args,repo)

		if err!= nil {
			return nil,err
		}
		
		fmt.Println(ver)
		fmt.Println(arch)

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








func resolveTargetArch(target core.Arch) core.Arch {

	if target == core.ArchUndefined{
		target = core.HostArch()
	}

	return target
}




func resolveArchVersion(args core.InstallArgs, repo core.Repo) (string, core.Arch, error) {



		target := resolveTargetArch(args.TargetArch)



		if args.Version.IsDefined() {
			return resolveRequestedVersion(args, repo, target)
		}


		return resolveLatestVersion(args, repo, target)
}



func resolveLatestVersion(args core.InstallArgs, repo core.Repo , resolvedArch core.Arch) (string, core.Arch, error) {

	v, ok, err := repo.GetLatestVersionForArch(args.PackageIdentifier, resolvedArch)
	if err != nil {
		return "", core.ArchUndefined, err
	}
	if ok {
		return v.Version, resolvedArch, nil
	}



	if core.HostArch() == core.ArchX64 {
		if args.TargetArch == core.ArchUndefined {
			v, ok, err := repo.GetLatestVersionForArch(args.PackageIdentifier, core.ArchX86)
			if err != nil {
				return "", core.ArchUndefined, err
			}
			if ok && v.ArchFallbackSafe {
				return v.Version, core.ArchX86, nil
			}

		}
	
	}

	return "", core.ArchUndefined,
    fmt.Errorf("package %s has no %q compatible version",
        args.PackageIdentifier,
        resolvedArch.String())

	
}



func resolveRequestedVersion(args core.InstallArgs, repo core.Repo , resolvedArch core.Arch) (string, core.Arch, error) {

	
	versions,err:= repo.GetVersions(args.PackageIdentifier)
	if err != nil {return "",core.ArchUndefined,err}

	if !slices.Contains(versions,args.Version.String()) {
		return "",core.ArchUndefined,fmt.Errorf("pkg %s does not have version %s",args.PackageIdentifier,args.Version.String())
	}

	// if the pkg no arch 

	versionInfo,err := repo.GetVersionInfo(args.PackageIdentifier,args.Version.String())
	if err != nil {return "",core.ArchUndefined,err}



	if slices.Contains(versionInfo.Archs, resolvedArch) {
		return args.Version.String(), resolvedArch, nil
	}

	if core.HostArch() == core.ArchX64 {
		if args.TargetArch == core.ArchUndefined && slices.Contains(versionInfo.Archs, core.ArchX86) && versionInfo.ArchFallbackSafe {
			return args.Version.String(), core.ArchX86, nil
		}
	}

	return "", core.ArchUndefined,
    fmt.Errorf("package %s version %s has no compatible architecture",
        args.PackageIdentifier,
        args.Version.String())

}


































/*
func mresolveArchVersion(args core.InstallArgs, repo core.Repo)(string,core.Arch,error){

		if args.Version.IsDefined() {
			//hanndle errors //fetch //should be sfve to asume not no pkgnotfound error
			versions,_ := repo.GetVersions(args.PackageIdentifier)

			if !slices.Contains(versions,args.Version.String()) {
				return "",core.ArchUndefined,fmt.Errorf("pkg %s does not have version %s",args.PackageIdentifier,args.Version.String())
			}

			// if the pkg no arch 

			versionInfo,err := repo.GetVersionInfo(args.PackageIdentifier,args.Version.String())
			if err != nil {return "",core.ArchUndefined,err}


		
			
			switch core.HostArch() {
				
			
				case core.ArchX64:

					if args.TargetArch == core.ArchUndefined {
						if slices.Contains(versionInfo.Archs,core.ArchX64){
							return args.Version.String(),core.ArchX64,nil
						}

						if archInfo.ArchFallbackSafe {
							if slices.Contains(versionInfo.Archs,core.ArchX86){
								return args.Version.String(),core.ArchX86,nil
							}
						}

						//error

					}

					if args.TargetArch == core.ArchX86 {
						if slices.Contains(versionInfo.Archs,core.ArchX86){
							return args.Version.String(),core.ArchX86,nil
						}

						// error 
					}


					// explicit x64/arm64 requests
					if slices.Contains(versionInfo.Archs,core.ArchX64){
						return args.Version.String(),core.ArchX64,nil
					}

					// error 



				case core.ArchX86:

					if slices.Contains(versionInfo.Archs,core.ArchX86){
							return args.Version.String(),core.ArchX86,nil
					}

					// error 

				default: // arm64
					if slices.Contains(versionInfo.Archs,core.ArchArm64){
						return args.Version.String(),core.ArchArm64,nil
					}

					// error 

			}

		}else{ // no defined verions

			if args.TargetArch == core.ArchUndefined {

				v,ok,err := repo.GetLatestVersionForArch(args.PackageIdentifier,core.HostArch()) // get latest verion for an arch 
				if err != nil {
					return "",core.ArchUndefined,err
				}


				if ok {
					return v.Version,core.HostArch(),nil
				}

				if core.HostArch() == core.ArchX64 {

					v,ok,_ := repo.GetLatestVersionForArch(args.PackageIdentifier,core.ArchX86)

					if ok {
						if v.ArchFallbackSafe {
							return v.Version,core.ArchX86,nil
						}
					}

				}

				// error 
			}


			switch core.HostArch() {

				case core.ArchX64:

					if args.TargetArch == core.ArchX64 {
						v,ok,err:= repo.GetLatestVersionForArch(args.PackageIdentifier,core.ArchX64)

						if err != nil {
							return "",core.ArchUndefined,err
						}
						if ok {
							return v.Version,core.ArchX64,nil
						}

						// error 
					}

					if args.TargetArch == core.ArchX86 {

						v,ok,err:= repo.GetLatestVersionForArch(args.PackageIdentifier,core.ArchX86)

						if err != nil {
							return "",core.ArchUndefined,err
						}
						if ok {
							return v.Version,core.ArchX86,nil
						}

						// error 
						
					}



				case core.ArchX86:
					v,ok,err := repo.GetLatestVersionForArch(args.PackageIdentifier,core.ArchX86)
					if err != nil {
						return "",core.ArchUndefined,err
					}

					if ok {
						return v.Version,core.ArchX86,nil
					}


					// error 

				default: //arm64
					v,ok,err := repo.GetLatestVersionForArch(args.PackageIdentifier,core.ArchArm64)

					if err != nil {
						return "",core.ArchUndefined,err
					}

					if ok {
						return v.Version,core.ArchArm64,nil
					}
	
					// error 
			}



		}


}






func cresolveArchVersion(args core.InstallArgs, repo core.Repo) (string, core.Arch, error) {





	if args.Version.IsDefined() {
		versions, err := repo.GetVersions(args.PackageIdentifier)
		if err != nil {
			return "", core.ArchUndefined, fmt.Errorf("fetching versions for %s: %w", args.PackageIdentifier, core.ErrFetch)
		}

		if !slices.Contains(versions, args.Version.String()) {
			return "", core.ArchUndefined, &core.VersionNotFoundError{
				Package: args.PackageIdentifier,
				Version: args.Version.String(),
			}
		}

		archInfo, err := repo.GetVersionInfo(args.PackageIdentifier, args.Version.String())
		if err != nil {
			return "", core.ArchUndefined, err
		}

		target := args.TargetArch
		if !target.IsDefined() {
			target = core.HostArch()
		}

		// exact match — always wins, no fallback logic needed
		if slices.Contains(archInfo.Archs, target) {
			return args.Version.String(), target, nil
		}

		// only fallback direction that's ever valid: target x64, x86 available, opted in
		if target == core.ArchX64 && slices.Contains(archInfo.Archs, core.ArchX86) && archInfo.ArchFallbackSafe {
			return args.Version.String(), core.ArchX86, nil
		}

		return "", core.ArchUndefined, &core.NoVersionForArchError{
			Package: args.PackageIdentifier,
			Version: args.Version.String(),
			Arch:    target,
		}
	}

	
	target := args.TargetArch
	if !target.IsDefined() {
		target = core.HostArch()
	}

	v, ok, err := repo.GetLatestVersionForArch(args.PackageIdentifier, target)
	if err != nil {
		return "", core.ArchUndefined, err
	}
	if ok {
		return v.Version, target, nil
	}

	// fallback only applies when target resolved to x64 (explicitly or via host default)
	if target == core.ArchX64 {
		v, ok, err := repo.GetLatestVersionForArch(args.PackageIdentifier, core.ArchX86)
		if err != nil {
			return "", core.ArchUndefined, err
		}
		if ok && v.ArchFallbackSafe {
			return v.Version, core.ArchX86, nil
		}
	}

	return "", core.ArchUndefined, &core.NoVersionForArchError{
		Package: args.PackageIdentifier,
		Arch:    target,
	}
}*/