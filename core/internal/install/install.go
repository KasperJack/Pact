package install

import (
	"fmt"
	//"go/version"
	//"slices"
	"errors"
	"github.com/kasperjack/pact/core"
	//"golang.org/x/tools/go/analysis/passes/nilfunc"
)




type manager struct {
    localState core.LocalState
    repo       core.Repo
    lockFile   core.LockFile
}



type installer struct {
    //args          core.InstallArgs
	
	manager manager

	metadata core.PackageInfo
    release  core.Release

    stagingDir    string
    installDir    string
    currentDir    string
}



func NewManaged(args core.InstallArgs, localState core.LocalState, repo core.Repo, lf core.LockFile) (*installer,error) {
	//resolve version 
	//resolve arch
	// resolve install type 

	i :=installer{}


	metadata, err := repo.LoadPackageInfo(args.PackageIdentifier)
	if err != nil {
		return nil,err
	}
	i.metadata = metadata


	
	r,err := resolveRelease(args,repo)

	if err != nil {
		return nil,err
	}


	fmt.Println(r)

	m := manager{
		localState: localState,
		repo: repo,
		lockFile: lf,
	}
	i. manager = m

	return &i,nil

}


func (i *installer) Run()error{
	// pcakge already exists at this point 

/*
	err := i.checkNotInstalled()
	if err != nil {
		return nil
	}

	

	i.checkVersionExists()

*/


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




func resolveRelease(args core.InstallArgs, repo core.Repo) (core.Release, error) {



		target := resolveTargetArch(args.TargetArch)



		if args.Version.IsDefined() {
			return resolveReleaseRequestedVersion(args, repo, target)
		}


		return resolveReleaseLatestVersion(args, repo, target)
}



func resolveReleaseLatestVersion(args core.InstallArgs, repo core.Repo , resolvedArch core.Arch) (core.Release, error) {

	fmt.Println("loading latest version")
	fullVer, err := repo.LoadLatestFullVersion(resolvedArch,args.PackageIdentifier)

	if err == nil {
		fmt.Printf("found latest version %s",fullVer)

		release, ferr := repo.LoadReleaseByFullVersion(resolvedArch,args.PackageIdentifier,fullVer)
	
		if ferr != nil {
			return core.Release{},ferr
		}
		return release,nil
	}

	if errors.Is(err, core.ErrPackageNotFoundForArch){

		if core.HostArch() == core.ArchX64 {
			if args.TargetArch == core.ArchUndefined {

				fFullVer, _ := repo.LoadLatestFullVersion(core.ArchX86,args.PackageIdentifier)

				release, aerr := repo.LoadReleaseByFullVersion(core.ArchX86,args.PackageIdentifier,fFullVer)
				fmt.Println(aerr)
				if aerr == nil {return release,nil }

			}
	
		}

	}



	return core.Release{},err

}










func resolveReleaseRequestedVersion(args core.InstallArgs, repo core.Repo, resolvedArch core.Arch) (core.Release, error) {

	

	release, err := repo.LoadReleaseByUpstreamVersion(resolvedArch, args.PackageIdentifier, args.Version.String())

	if err == nil {
		return release, nil
	}

	if errors.Is(err, core.ErrPackageNotFoundForArch) {
		if core.HostArch() == core.ArchX64 {
			if args.TargetArch == core.ArchUndefined {

				release, fallbackErr := repo.LoadReleaseByUpstreamVersion(core.ArchX86, args.PackageIdentifier, args.Version.String())
				if fallbackErr == nil {
					return release, nil
				}
			}
		}
	}


	// version not found /io
	return core.Release{}, err
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