package main

import(
	
)


func build(path string) error{

	raw,err := loadFile(path)

	if err != nil {
		return err
	}

	err = parseConfig(raw)

	if err != nil {
		return err
	}


	return nil
}
