package connectdb

import (
	"fmt"
	"techtrain/techdb"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// update db
func ConnUpdateName(x string, n string) {
	// try to connect db
	dsn := rUsername + ":" + rPassword + "@" + rProtocol + "(" + rAddress + ")" + "/" + rDbname
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		panic(err)
	}

	//insert xtoken and name
	var user techdb.User
	SQLrequest := db.Model(&user).Where("Xtoken = ?", x).Update("Name", "Alice")
	err = SQLrequest.Error
	if err != nil {
		fmt.Println(err)
	}

	return
}
