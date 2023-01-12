package limitedgacha

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Characterinfo struct {
	Characterid uint `gorm:"primaryKey"`
	Name        string
	Stdpower    uint
}

// read db
func ConnSetInfo() {
	// try to connect db
	dsn := Username + ":" + Password + "@" + Protocol + "(" + Address + ")" + "/" + Dbname

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		panic(err)
	}

	var characterinfo = []Characterinfo{
		{Characterid: 11, Name: "Skadi the Corrupting Heart", Stdpower: 3000},
	}

	//insert new character information
	SQLrequest := db.Create(characterinfo)
	err = SQLrequest.Error
	if err != nil {
		fmt.Println(err)
	}

}
