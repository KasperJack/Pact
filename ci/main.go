package main

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"

	"log"
	"github.com/kasperjack/pact/core/parce"
	"os"

)






func main() {


	mfsSrc, err := os.ReadFile("mfs.hcl")
	if err != nil {
		log.Fatalf("reading mfs.hcl: %v", err)
	}

	m, diags := parce.Manifest(mfsSrc)
	if diags.HasErrors() {
		printDiags(diags)
		log.Fatal("parse failed")
	}

	vm, diags := ValidateManifest(m)
	if diags.HasErrors() {
		printDiags(diags)
		log.Fatal("validation failed")
	}


	
	interSrc, err := os.ReadFile("inter.hcl")
	if err != nil {
		log.Fatalf("reading inter.hcl: %v", err)
	}

	inter, diags := parce.Interface(interSrc)
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

