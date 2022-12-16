package techdb

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// read db
func ConnSetInfo() {
	// try to connect db
	dsn := rUsername + ":" + rPassword + "@" + rProtocol + "(" + rAddress + ")" + "/" + rDbname

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		panic(err)
	}

	var characterinfo = []Characterinfo{
		{Characterid: 1, Name: "Skadi", Stdpower: 2000},
		{Characterid: 2, Name: "Sliverash", Stdpower: 2000},
		{Characterid: 3, Name: "CH'EN", Stdpower: 2000},
		{Characterid: 4, Name: "Skyfire", Stdpower: 1500},
		{Characterid: 5, Name: "Sora", Stdpower: 1500},
		{Characterid: 6, Name: "Astesia", Stdpower: 1500},
		{Characterid: 7, Name: "ShiraYuki", Stdpower: 1000},
		{Characterid: 8, Name: "Vigna", Stdpower: 1000},
		{Characterid: 9, Name: "Frostleaf", Stdpower: 1000},
		{Characterid: 10, Name: "Jaye", Stdpower: 1000},
	}

	//insert new character information
	SQLrequest := db.Create(characterinfo)
	err = SQLrequest.Error
	if err != nil {
		fmt.Println(err)
	}

}
