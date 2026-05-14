package main

import(
	
)


func build(path string) error{

	raw,err := loadFile(path)

	if err != nil {
		return err
	}

	p,err := parseConfig(raw)

	if err != nil {
		return err
	}

	p.Validate()



	return nil
}
