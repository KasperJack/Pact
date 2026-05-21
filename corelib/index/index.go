package index

import (
	//"gorm.io/driver/sqlite"
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	//"github.com/BurntSushi/toml"
	//"log"
	"Pact/corelib/model"
	"fmt"
)



type User struct {
    ID   uint
    Name string
    Age  int
}

func Rebuild() error {

	var pacakges []model.PackageDB
/*
	var pk model.PackageT

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
*/

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        panic(err)
    }

	err = db.AutoMigrate(&model.PackageDB{},&model.ReleaseDB{})
	if err != nil {
		log.Fatal(err)
	}


	pkg := model.PackageDB{
    PackageIdentifier: "react",
    Name:              "React",
    Description:       "A JavaScript library",
    Versioning:        "semver",
    Releases: []model.ReleaseDB{
        {
            Version: "20.0.0",
            Hash:    "abc123",
            URL:     "https://...",
        },
        {
            Version: "20.0.0",
            Hash:    "def456",
            URL:     "https://...",
        },
    },
}

	result := db.Create(&pkg)
	if result.Error != nil {
		log.Fatal(result.Error) // will catch duplicate PK violations
	}


	//db.Create(&pk) // error ? 
	return nil
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