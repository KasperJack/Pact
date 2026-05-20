package index

import (
    //"gorm.io/driver/sqlite"
    "gorm.io/gorm"
	"github.com/glebarez/sqlite"
	
	"github.com/BurntSushi/toml"
	"log"
	"fmt"
	"Pact/corelib/model"
)



type User struct {
    ID   uint
    Name string
    Age  int
}

func LoadToml() {

	var pk model.Package

	_, err := toml.DecodeFile("C:\\Users\\Aya\\Desktop\\pact-tools\\test-buckets\\defult\\as\\asshoe\\package.toml", &pk)
	if err != nil {
		log.Fatal(err)
	}

	err = pk.ValidateOnRead()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(pk.Name)
	fmt.Println(pk.Versioning)
	fmt.Println(pk.Description)
	fmt.Println(pk.PackageIdentifier)


	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        panic(err)
    }

	db.AutoMigrate(&model.Package{})
	db.Create(&pk) // error ? 
}



func Test() {
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        panic(err)
    }

    db.AutoMigrate(&User{})

    db.Create(&User{
        Name: "Kasper",
        Age:  22,
    })

    // Query
    var user User
    db.First(&user)

    fmt.Println(user.Name)
}