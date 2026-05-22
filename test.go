package pacttools
/*func NewPackage(luaData []byte, client *client.RepoClient) (pipeline.PublisherPackage,error) {

	L := lua.NewState()
	defer L.Close()

	var pacakge model.Package

	L.SetGlobal("package", L.NewFunction(func(L *lua.LState) int {

		tbl := L.CheckTable(1)

	
		getRequiredString := func(key string) string {

			val := tbl.RawGetString(key)

			if val.Type() == lua.LTNil {

                L.RaiseError("missing required field: %s", key)
                
            }

			if val.Type() != lua.LTString {
                L.RaiseError("field '%s' must be a string", key)
                
            }

			str := val.String()

			if str == "" {
                L.RaiseError("field '%s' cannot be empty", key)
               
            }
			return str
		}





		pacakge = model.Package{
			PackageIdentifier: getRequiredString("package_identifier"),
			Name:              tbl.RawGetString("name").String(),
			Versioning:        tbl.RawGetString("versioning").String(),
			Description:       tbl.RawGetString("description").String(),
			Homepage:          tbl.RawGetString("homepage").String(),
			License:           tbl.RawGetString("license").String(),
		}

		return 0
	}))

	// run lua file
	if err := L.DoString(string(luaData)); err != nil {
		return nil,err
	}





	fmt.Println(pacakge.PackageIdentifier)
	fmt.Println(pacakge.Name)
	fmt.Println(pacakge.Versioning)

	fmt.Println(pacakge.Description)
	fmt.Println(pacakge.Homepage)
	fmt.Println(pacakge.License)

	return &PackageContext{Model: &pacakge, Client: client},nil
}*/