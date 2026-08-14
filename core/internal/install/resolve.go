package install




import (
	"fmt"
	"slices"
	"errors"
	"github.com/kasperjack/pact/core"
	
)










func resolveTargetArch(target core.Arch) core.Arch {

	if target == core.ArchUndefined{
		target = core.HostArch()
	}

	return target
}





func resolveRelease(args core.InstallArgs, repo core.Repo) (core.Release, error) {



		//target := resolveTargetArch(args.TargetArch)



		if args.Version.IsDefined() {
			return resolveReleaseRequestedVersion(args, repo,)
		}


		return resolveReleaseLatestVersion(args, repo)
}








func resolveReleaseLatestVersion(args core.InstallArgs, repo core.Repo) (core.Release, error) {


	if args.TargetArch == core.ArchUndefined {

		host := core.HostArch()
		compat := host.CompatibleArchs()

		for _,arch := range compat {

			p, err := repo.Package(arch,args.PackageIdentifier)

			if err != nil {
				if errors.Is(err, core.ErrPkgNotFound) {
					continue
				}
				return core.Release{}, err
			}

			latestVersion := p.LatestReleaseVersion()

			return p.Release(latestVersion)

		}

		return core.Release{}, fmt.Errorf("package %s not found for any compatible architecture ([debug]:tried: %v)", args.PackageIdentifier, compat)

	}

	// explicit target arch specified, try to get the package for that arch


	p, err := repo.Package(args.TargetArch,args.PackageIdentifier)

	if err != nil {

		if errors.Is(err, core.ErrPkgNotFound){

			return core.Release{}, fmt.Errorf("package %s not found for architecture %s", args.PackageIdentifier, args.TargetArch.String())
		}

		return core.Release{}, err	
	}

	latestVersion := p.LatestReleaseVersion()

	return p.Release(latestVersion)



}






func resolveReleaseRequestedVersion(args core.InstallArgs, repo core.Repo) (core.Release, error) {


	if args.TargetArch == core.ArchUndefined {

		host := core.HostArch()
		compat := host.CompatibleArchs()

		// make a slice of existing packages for the compatible archs

		pks := []core.Package{}

		for _,arch := range compat {

			p, err := repo.Package(arch,args.PackageIdentifier)
		

			if err != nil {
				if errors.Is(err, core.ErrPkgNotFound) {
					continue
				}

				return core.Release{}, err
			}

			pks = append(pks, p)
		}

		if len(pks) == 0 {
			return core.Release{}, fmt.Errorf("package %s not found for any compatible architecture ([debug]:tried: %v)", args.PackageIdentifier, compat)
		}



		for _,p := range pks {

			if slices.Contains(p.UpstreamVersions(),args.Version.String()){

				return p.Release(p.LatestReleaseForUpstream(args.Version.String()))
			}
			

		}



		return core.Release{}, fmt.Errorf("version %s of package %s not found for any compatible architecture ([debug]:tried: %v)", args.Version.String(), args.PackageIdentifier, compat)		

	}




	p, err := repo.Package(args.TargetArch,args.PackageIdentifier)

	if err != nil {  // fallback here 

		if errors.Is(err, core.ErrPkgNotFound) {

			return core.Release{}, fmt.Errorf("package %s not found for architecture %s", args.PackageIdentifier, args.TargetArch.String())
		}

		return core.Release{}, err
	}


	if !slices.Contains(p.UpstreamVersions(),args.Version.String()){

		return core.Release{}, fmt.Errorf("version %s of package %s not found for architecture %s", args.Version.String(), args.PackageIdentifier, args.TargetArch.String())

	}

	return p.Release(p.LatestReleaseForUpstream(args.Version.String())) 


}