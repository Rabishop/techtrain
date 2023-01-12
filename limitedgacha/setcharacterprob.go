package limitedgacha

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Characterprobwithlimit struct {
	Listid      uint `gorm:"primaryKey"`
	Characterid uint `gorm:"primaryKey"`
	Prob        uint
	Number      uint
}

// read db
func ConnSetProb() {
	// try to connect db
	dsn := Username + ":" + Password + "@" + Protocol + "(" + Address + ")" + "/" + Dbname
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		panic(err)
	}

	var Characterprobwithlimit = []Characterprobwithlimit{
		{Listid: 1, Characterid: 1, Prob: 4, Number: 10},
		{Listid: 1, Characterid: 2, Prob: 4, Number: 10},
		{Listid: 1, Characterid: 3, Prob: 4, Number: 10},
		{Listid: 1, Characterid: 4, Prob: 1000, Number: 1000},
		{Listid: 1, Characterid: 5, Prob: 1000, Number: 1000},
		{Listid: 1, Characterid: 6, Prob: 1000, Number: 1000},
		{Listid: 1, Characterid: 7, Prob: 23997, Number: 999999999},
		{Listid: 1, Characterid: 8, Prob: 23997, Number: 999999999},
		{Listid: 1, Characterid: 9, Prob: 23997, Number: 999999999},
		{Listid: 1, Characterid: 10, Prob: 23997, Number: 999999999},
		{Listid: 1, Characterid: 11, Prob: 1000, Number: 1000},

		{Listid: 2, Characterid: 1, Prob: 4, Number: 1},
		{Listid: 2, Characterid: 2, Prob: 4, Number: 1},
		{Listid: 2, Characterid: 3, Prob: 4, Number: 1},
		{Listid: 2, Characterid: 4, Prob: 1000, Number: 100},
		{Listid: 2, Characterid: 5, Prob: 1000, Number: 100},
		{Listid: 2, Characterid: 6, Prob: 1000, Number: 100},
		{Listid: 2, Characterid: 7, Prob: 23997, Number: 999999999},
		{Listid: 2, Characterid: 8, Prob: 23997, Number: 999999999},
		{Listid: 2, Characterid: 9, Prob: 23997, Number: 999999999},
		{Listid: 2, Characterid: 10, Prob: 23997, Number: 999999999},
		{Listid: 2, Characterid: 11, Prob: 1000, Number: 100},
	}

	//insert new character probability
	SQLrequest := db.Save(Characterprobwithlimit)
	err = SQLrequest.Error
	if err != nil {
		fmt.Println(err)
	}

}
