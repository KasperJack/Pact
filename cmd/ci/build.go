package main

import(
	
)


func build(path string) error{
	
	raw,err := loadFile(path)

	if err != nil {
		return err
	}



	err = parseluaConfig(raw)

	if err != nil {
		return err
	}



	return nil
}
