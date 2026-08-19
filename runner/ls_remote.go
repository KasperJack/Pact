package main

import (
	//"fmt"
	"strings"
	//"errors"
	"fmt"
	//"os"
	//"path/filepath"

	//"github.com/go-delve/delve/pkg/version"
	"github.com/kasperjack/pact/core"
	//"github.com/kasperjack/pact/core/repo"
)



type archVersions struct {
	arch core.Arch
	versions []string

}

func LsRemote(packageIdentifier string) error {


	/*

	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exeDir := filepath.Dir(exePath)


	r,err := repo.NewLocalRepo(filepath.Join(exeDir, "repo"))  
	if err != nil {
		return err
	}
	
	host := core.HostArch()


	switch host {


	case core.ArchX64:
		var avs []archVersions

		p64 , err := r.Package(core.ArchX64,packageIdentifier)

		if err != nil {

			if !errors.Is(err, core.ErrPkgNotFound) {
				return err
			}
		

		}else{ // nil error 

			avs = append(avs, archVersions{
			arch:     core.ArchX64,
			versions: p64.UpstreamVersions(),
			})

		}


		p86 , err := r.Package(core.ArchX86,packageIdentifier)

		if err != nil {

			if !errors.Is(err, core.ErrPkgNotFound) {
				return err
			}

		}else{

			avs = append(avs, archVersions{
			arch:     core.ArchX86,
			versions: p86.UpstreamVersions(),
			})

		}


		if len(avs) == 0 {
			return fmt.Errorf("pkg not found")
		}

		listVersions(avs)


	default: // x86 //arm64

		p , err := r.Package(host,packageIdentifier)

		if err != nil {return  err} // error includes pkg not found

		listVersions([]archVersions{{
			arch:     host,
			versions: p.UpstreamVersions(),
		}})
			

	}



	return nil	
	*/

	return fmt.Errorf("ls-remote not implmented yet")
}





func listVersions(avs []archVersions) {

		// do print here for now, version order don't matter for now 
	// use arch.string() to get the cli reprsantation 
	
	/// single arch 
	// 1.2.3
	// 1.2.6
	// 1.6.3

	///multi arch

	// 1.2.3 [x64][x86] 
	// 1.2.7 [x64][x86]
	// 1.2.5 [x64][x86]
	// 1.1.3 [x86]
	// 1.3.3 [x64] 





	if len(avs) == 1 {
		for _, v := range avs[0].versions {
			fmt.Println(v)
		}
		return
	}

	// multi-arch: merge into a union, tracking which archs have each version
	archsByVersion := make(map[string][]core.Arch)
	var order []string // preserves first-seen order across archs

	for _, av := range avs {
		for _, v := range av.versions {
			if _, seen := archsByVersion[v]; !seen {
				order = append(order, v)
			}
			archsByVersion[v] = append(archsByVersion[v], av.arch)
		}
	}

	for _, v := range order {
		var tags strings.Builder
		for _, a := range archsByVersion[v] {
			tags.WriteString("[")
			tags.WriteString(a.String())
			tags.WriteString("]")
		}
		fmt.Printf("%s %s\n", v, tags.String())
	}
}