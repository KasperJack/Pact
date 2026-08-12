package install

import (
	"fmt"
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

	
	_,ok := lf.GetInstalled(args.PackageIdentifier)
	if ok {
		return nil,fmt.Errorf("package %s is already installed",args.PackageIdentifier)
	}
	

	//lf.Test()
	

	i :=installer{}

	md,err  := repo.LoadPackageInfo(args.PackageIdentifier)
	if err != nil {
		return nil,err
	}

	i.metadata = md


	r,err := resolveRelease(args,repo)
	if err != nil {
		return nil,err
	}


	i.release = r

	m := manager{
		localState: localState,
		repo: repo,
		lockFile: lf,
	}
	i. manager = m

	return &i,nil

}

// Run executes the installation process
func (i *installer) Run()error{
	
	//fmt.Println("installing:")
	//fmt.Println(i.metadata.Package)
	//fmt.Println(i.release.Architecture)
	//fmt.Println(i.release.URL)
	//fmt.Println(i.release.FullVersion)
	//fmt.Println(i.release.UpstreamVersion)







	return nil
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