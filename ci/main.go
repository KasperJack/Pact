package main

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"

	"log"
	"github.com/kasperjack/pact/core/parce"
	"os"
	"path/filepath"

)






func main() {

	pkgPath := "C:/Users/kasper/Documents/projects/pact/bin/repo/packages/windirstat"

	archRelease := filepath.Join(pkgPath, "2.7.0", "x64", "1")



	packageSrc, err := os.ReadFile(filepath.Join(pkgPath, "package.hcl"))

	if err != nil {
		log.Fatalf("reading package.hcl: %v", err)
	}
	_, diags := parce.PackageInfo(packageSrc)

	if diags.HasErrors() {
		printDiags(diags)
		log.Fatal("package info parse failed")
	}




	manifestSrc, err := os.ReadFile(filepath.Join(archRelease, "manifest.hcl"))
	if err != nil {
		log.Fatalf("reading manifest.hcl: %v", err)
	}

	m, diags := parce.Manifest(manifestSrc)
	if diags.HasErrors() {
		printDiags(diags)
		log.Fatal("parse failed")
	}

	vm, diags := ValidateManifest(m)
	if diags.HasErrors() {
		printDiags(diags)
		log.Fatal("validation failed")
	}





	
	interfaceSrc, err := os.ReadFile(filepath.Join(archRelease, "interface.hcl"))
	if err != nil {
		log.Fatalf("reading interface.hcl: %v", err)
	}

	inter, diags := parce.Interface(interfaceSrc)
	if diags.HasErrors() {
		printDiags(diags)
		log.Fatal("interface parse failed")
	}

	_ = vm
	_ = inter
}

func printDiags(diags hcl.Diagnostics) {
	for _, d := range diags {
		fmt.Println(d.Error())
	}
}

