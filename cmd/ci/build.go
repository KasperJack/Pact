package main

import(
	"Pact/corelib/stage"
)


func build(path string) error{

	raw,err := stage.LoadFile(path)

	if err != nil {
		return err
	}

	_,err = stage.ParseConfig(raw)

	if err != nil {
		return err
	}






	return nil
}
