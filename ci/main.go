package main

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"

	"log"
	"os"
	"path/filepath"

	"github.com/kasperjack/pact/core"
	"github.com/kasperjack/pact/core/parce"
)






func main() {

	pkgPath := "C:/Users/kasper/Documents/projects/pact/bin/repo/packages/windirstat"

	archRelease := filepath.Join(pkgPath, "2.7.0", "x64", "1")



	packageSrc, err := os.ReadFile(filepath.Join(pkgPath, "package.hcl"))

	if err != nil {
		log.Fatalf("reading package.hcl: %v", err)
	}


	pkg, diags := parce.PackageInfo(packageSrc)

	if diags.HasErrors() {
		printDiags(diags)
		log.Fatal("package info parse failed")
	}

	
	/*
	
	interfaceSrc, err := os.ReadFile(filepath.Join(archRelease, "interface.hcl"))
	if err != nil {
		log.Fatalf("reading interface.hcl: %v", err)
	}

	inter, diags := parce.Interface(interfaceSrc,pkg.Scopes)
	if diags.HasErrors() {
		printDiags(diags)
		log.Fatal("interface parse failed")
	}
	*/

	fmt.Println(pkg.Scopes)


	manifestSrc, err := os.ReadFile(filepath.Join(archRelease, "manifest.hcl"))
	if err != nil {
		log.Fatalf("reading manifest.hcl: %v", err)
	}

	m, diags := parce.Manifest(manifestSrc,[]core.Scope{core.ScopeUser,core.ScopeSystem})
	if diags.HasErrors() {
		printDiags(diags)
		log.Fatal("parse failed")
	}

	s,ok  :=m.Scope[core.ScopeUser]

	if !ok {
		log.Fatal("user scope not found")
	}

	fmt.Println(len(s.Blocks))

	var v core.Visitor = core.DryRunner{}

	for _, b := range s.Blocks {
		if err := b.Accept(v); err != nil {
			log.Fatal(err)
		}
	}

}

func printDiags(diags hcl.Diagnostics) {
	for _, d := range diags {
		fmt.Println(d.Error())
	}
}

