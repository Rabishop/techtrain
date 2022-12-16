package gacha

import (
	"fmt"
	"techtrain/techdb"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func ConnReadInfo(x string, userinventory *[]techdb.Userinventory) int {
	// try to connect db
	dsn := rUsername + ":" + rPassword + "@" + rProtocol + "(" + rAddress + ")" + "/" + rDbname
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		fmt.Println(err)
		return 300
	}

	// get userid
	var user techdb.User
	SQLrequest := db.Where("Xtoken = ?", x).First(&user)
	err = SQLrequest.Error
	if err != nil {
		fmt.Println(err)
		return 301
	}

	// select user
	SQLrequest = db.Where("Userid = ?", user.Userid).Find(&userinventory)
	err = SQLrequest.Error
	if err != nil {
		fmt.Println(err)
		return 301
	}
	return 100
}
