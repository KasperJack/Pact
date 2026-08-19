package install




import (
	"fmt"
	"slices"
	//"errors"
	"github.com/kasperjack/pact/core"
	
)











																// retrun also metadata 
func resolveArchRelease(args core.InstallArgs, repo core.Repo) (core.PackageInfo, core.Release, error) {


	
	md,err  := repo.LoadPackageInfo(args.PackageIdentifier)

	if err != nil {

		return core.PackageInfo{},core.Release{},err
	}



	if args.TargetArch == core.ArchUndefined {

		foundArch := false

		for _,arch := range core.HostArch().CompatibleArchs() {
			if slices.Contains(md.Architectures,arch) {
				foundArch = true
				break
			}
			
		}

		if !foundArch {
			return core.PackageInfo{},core.Release{} ,fmt.Errorf("package %s does not support any compatible architecture for this host", args.PackageIdentifier)
		}

	}else{

		if !slices.Contains(md.Architectures,args.TargetArch) {
			return core.PackageInfo{},core.Release{} ,fmt.Errorf("package %s does not support architecture %s", args.PackageIdentifier, args.TargetArch.String())
		}

	}









		if args.Version.IsDefined() {
			r , err := resolveReleaseRequestedVersion(args, repo)

			if err != nil {
				return core.PackageInfo{},core.Release{},err
			}

			
			return md,r,nil
		}



		r, err := resolveReleaseLatestVersion(args, repo)

		if err != nil {
				return core.PackageInfo{},core.Release{},err
			}

		return md,r,nil
		
}








func resolveReleaseLatestVersion(args core.InstallArgs, repo core.Repo) (core.Release, error) {


	index, err := repo.LoadPackageIndex(args.PackageIdentifier)

	if err != nil {
		return core.Release{}, err
	}


	archMap := index.ArchMap.AsMap()


	if args.TargetArch == core.ArchUndefined {

		host := core.HostArch()
		compat := host.CompatibleArchs()

		for _,arch := range compat {

			


			archInfo,ok := archMap[arch]

			if ok{

				return repo.LoadArchRelease(args.PackageIdentifier,archInfo.Version,arch,archInfo.Revision)

			}

		}



		return core.Release{}, fmt.Errorf("package %s has no avalable relase ? for any compatible architecture ", args.PackageIdentifier)

	}




	archInfo,ok := archMap[args.TargetArch]

	if ok{

		return repo.LoadArchRelease(args.PackageIdentifier,archInfo.Version,args.TargetArch,archInfo.Revision)

	}

	return core.Release{}, fmt.Errorf("package %s has no avalable relase ? for %s architecture ", args.PackageIdentifier,args.TargetArch.String())
}













func resolveReleaseRequestedVersion(args core.InstallArgs, repo core.Repo) (core.Release, error) {


	index, err := repo.LoadPackageIndex(args.PackageIdentifier)

	if err != nil {
		return core.Release{}, err
	}

	
	if !slices.Contains(index.Versions,args.Version.String()) {
		return core.Release{}, fmt.Errorf("version %s of package %s not found", args.Version.String(), args.PackageIdentifier)
	}




	if args.TargetArch == core.ArchUndefined {

		host := core.HostArch()
		compat := host.CompatibleArchs()


		for _,arch := range compat {

			s, err := repo.LoadArchStatus(args.PackageIdentifier,args.Version.String(),arch)

			if err != nil {
				return core.Release{}, err
			}


			if s.Status == "available" {
				return repo.LoadArchRelease(args.PackageIdentifier,args.Version.String(),arch,s.CurrentRevision)
			}



		}


		return core.Release{}, fmt.Errorf("package %s version %s not available for any compatible architecture ([debug]:tried: %v)", args.PackageIdentifier, args.Version.String(), compat)

	}


	s, err := repo.LoadArchStatus(args.PackageIdentifier,args.Version.String(),args.TargetArch)

	if err != nil {
			return core.Release{}, err
		}

	if s.Status == "available" {
		return repo.LoadArchRelease(args.PackageIdentifier,args.Version.String(),args.TargetArch,s.CurrentRevision)
	}



	return core.Release{}, fmt.Errorf("package %s version %s not available for architecture %s", args.PackageIdentifier, args.Version.String(), args.TargetArch.String())


}