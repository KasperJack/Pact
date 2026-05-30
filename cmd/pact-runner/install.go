package main
import "os"

func install (pkg string) error {


	info, err := os.Stat(pkg)
	if err != nil || !info.IsDir() {
		return err
	}




	

	return nil
}