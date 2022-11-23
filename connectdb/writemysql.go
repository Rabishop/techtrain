package connectdb

import (
	"fmt"
	"techtrain/techdb"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// write db
func ConnWriteName(x string, n string) int {
	// try to connect db
	dsn := rUsername + ":" + rPassword + "@" + rProtocol + "(" + rAddress + ")" + "/" + rDbname
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		panic(err)
	}

	//insert xtoken and name
	var user techdb.User
	user.Name = n
	user.Xtoken = x
	SQLrequest := db.Create(&user)
	err = SQLrequest.Error
	if err != nil {
		fmt.Println(err)
		return 301
	}

	return 100
}
