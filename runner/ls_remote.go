package main

import (
	//"fmt"
	"os"
	"path/filepath"

	"github.com/kasperjack/pact/core"
	"github.com/kasperjack/pact/core/repo"
)

func LsRemote(packageIdentifier string) error {

	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exeDir := filepath.Dir(exePath)


	r,err := repo.NewLocalRepo(filepath.Join(exeDir, "repo"))  // test-buckets
	if err != nil {
		return err
	}
	



	switch core.HostArch(){

	case core.ArchX64:

	case core.ArchX86:	
		r.HasPackage(core.ArchX86,packageIdentifier)


	case core.ArchArm64:

	
	}



	return nil	
}




func listSingleArch(Target core.Arch, repo core.Repo)error{

	return nil
}


func listMultiArch(Target core.Arch, repo core.Repo)error{

	return nil
}
// add unexpected errors for io and readding files 