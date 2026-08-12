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



		target := resolveTargetArch(args.TargetArch)



		if args.Version.IsDefined() {
			return resolveReleaseRequestedVersion(args, repo, target)
		}


		return resolveReleaseLatestVersion(args, repo, target)
}








func resolveReleaseLatestVersion(args core.InstallArgs, repo core.Repo , resolvedArch core.Arch) (core.Release, error) {



	p, err := repo.Package(resolvedArch,args.PackageIdentifier)

	if err != nil {

		if errors.Is(err, core.ErrPkgNotFound) && core.HostArch() == core.ArchX64 && args.TargetArch == core.ArchUndefined {

			fp, ferr := repo.Package(core.ArchX86,args.PackageIdentifier)
			if ferr != nil {return core.Release{},ferr}

			return fp.Release(fp.LatestReleaseVersion())

		}



		return core.Release{}, err	
	}

	latestVersion := p.LatestReleaseVersion()

	return p.Release(latestVersion)



}



func resolveReleaseRequestedVersion(args core.InstallArgs, repo core.Repo, resolvedArch core.Arch) (core.Release, error) {





	p, err := repo.Package(resolvedArch,args.PackageIdentifier)

	if err != nil {  // fallback here 

		if errors.Is(err, core.ErrPkgNotFound) && core.HostArch() == core.ArchX64 && args.TargetArch == core.ArchUndefined {

	
			fp, ferr := repo.Package(core.ArchX86,args.PackageIdentifier)

			if ferr != nil {return core.Release{},ferr}

			if !slices.Contains(fp.UpstreamVersions(),args.Version.String()){
				return core.Release{},fmt.Errorf("version %s not found",args.Version.String()) // did not find a pkg for x64 found x86 no requested version
			}

			return fp.Release(fp.LatestReleaseForUpstream(args.Version.String())) 



		}



		return core.Release{}, err
	}


	if !slices.Contains(p.UpstreamVersions(),args.Version.String()){

		if core.HostArch() == core.ArchX64 && args.TargetArch == core.ArchUndefined {


			fp, ferr := repo.Package(core.ArchX86,args.PackageIdentifier)

			if ferr != nil {return core.Release{},ferr}

			if !slices.Contains(fp.UpstreamVersions(),args.Version.String()){
				return core.Release{},fmt.Errorf("version %s not found",args.Version.String()) // found x64 pkg no version found x86 no version 
			}

			return fp.Release(fp.LatestReleaseForUpstream(args.Version.String())) 


		}


		return core.Release{},fmt.Errorf("version %s not found",args.Version.String())
	}

	return p.Release(p.LatestReleaseForUpstream(args.Version.String())) 


}