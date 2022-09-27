package techdb

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// read db
func ConnSetProb() {
	// try to connect db
	dsn := rUsername + ":" + rPassword + "@" + rProtocol + "(" + rAddress + ")" + "/" + rDbname
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		panic(err)
	}

	var characterprob = []Characterprob{
		{Listid: 1, Characterid: 1, Prob: 4},
		{Listid: 1, Characterid: 2, Prob: 4},
		{Listid: 1, Characterid: 3, Prob: 4},
		{Listid: 1, Characterid: 4, Prob: 1000},
		{Listid: 1, Characterid: 5, Prob: 1000},
		{Listid: 1, Characterid: 6, Prob: 1000},
		{Listid: 1, Characterid: 7, Prob: 24247},
		{Listid: 1, Characterid: 8, Prob: 24247},
		{Listid: 1, Characterid: 9, Prob: 24247},
		{Listid: 1, Characterid: 10, Prob: 24247},

		{Listid: 2, Characterid: 1, Prob: 8},
		{Listid: 2, Characterid: 4, Prob: 1000},
		{Listid: 2, Characterid: 5, Prob: 1000},
		{Listid: 2, Characterid: 7, Prob: 24498},
		{Listid: 2, Characterid: 8, Prob: 24498},
		{Listid: 2, Characterid: 9, Prob: 24498},
		{Listid: 2, Characterid: 10, Prob: 24498},

		{Listid: 3, Characterid: 2, Prob: 8},
		{Listid: 3, Characterid: 4, Prob: 1000},
		{Listid: 3, Characterid: 6, Prob: 1000},
		{Listid: 3, Characterid: 7, Prob: 24498},
		{Listid: 3, Characterid: 8, Prob: 24498},
		{Listid: 3, Characterid: 9, Prob: 24498},
		{Listid: 3, Characterid: 10, Prob: 24498},

		{Listid: 4, Characterid: 3, Prob: 8},
		{Listid: 4, Characterid: 5, Prob: 1000},
		{Listid: 4, Characterid: 6, Prob: 1000},
		{Listid: 4, Characterid: 7, Prob: 24498},
		{Listid: 4, Characterid: 8, Prob: 24498},
		{Listid: 4, Characterid: 9, Prob: 24498},
		{Listid: 4, Characterid: 10, Prob: 24498},
	}

	//insert new character probability
	SQLrequest := db.Create(characterprob)
	err = SQLrequest.Error
	if err != nil {
		fmt.Println(err)
	}

}
