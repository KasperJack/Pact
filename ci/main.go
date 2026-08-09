
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	//"github.com/knqyf263/go-deb-version"
	//"github.com/kasperjack/pact/core"
	"github.com/kasperjack/pact/core"
	"github.com/kasperjack/pact/core/repo"
)




func main(){


	if len(os.Args) > 2 {
		
		switch os.Args[1] {

		case "build":

			buildCmd(os.Args[2])

			
		}
	
	}


	fmt.Fprintln(os.Stderr,"use build <pkg>")
	os.Exit(1)
}










func buildCmd(pkg string) {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("error: failed to get executable path: %v", err)
	}

	exeDir := filepath.Dir(exePath)
	log.Printf("using repo in %s", exeDir)

	r, err := repo.NewLocalRepo(filepath.Join(exeDir, "repo"))
	if err != nil {
		log.Fatalf("error: failed to initialize repository at %s: %v", exeDir, err)
	}

	ok, err := r.PackageExists(pkg)
	if err != nil {
		log.Fatalf("error: failed checking package existence for %s: %v", pkg, err)
	}

	if !ok {
		log.Printf("pkg %q not found in repository index", pkg)

		// Check if release metadata exists for specific architectures even if package index missed it
		checkRelease(core.ArchX64, pkg, r)
		checkRelease(core.ArchX86, pkg, r)
		checkRelease(core.ArchArm64, pkg, r)

		os.Exit(1)
	}

	p, err := r.LoadPackageInfo(pkg)
	if err != nil {
		log.Fatalf("error: failed loading package info for %s: %v", pkg, err)
	}

	fmt.Println("Package:", p.Package)
	fmt.Println("Name:", p.Name)
	fmt.Println("License:", p.License)
	fmt.Println("Homepage:", p.Homepage)
	fmt.Println("Description:", p.Description)

	// Check release availability across target architectures
	checkRelease(core.ArchX64, pkg, r)
	checkRelease(core.ArchX86, pkg, r)
	checkRelease(core.ArchArm64, pkg, r)

	os.Exit(0)
}

func checkRelease(arch core.Arch, pkg string, r core.Repo) {
	ok, err := r.PackageExistsForArch(arch, pkg)
	if err != nil {
		log.Printf("error: failed checking release for arch %s: %v", arch, err)
		return
	}

	if ok {
		log.Printf("found release for arch: %s", arch)
		return
	}

	log.Printf("no release found for arch: %s", arch)
}

