package connectdb

import (
	"fmt"
	"techtrain/techdb"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func ConnReadName(x string) string {
	// try to connect db
	dsn := rUsername + ":" + rPassword + "@" + rProtocol + "(" + rAddress + ")" + "/" + rDbname
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		panic(err)
	}

	//select all
	var user techdb.User
	SQLrequest := db.Where("Xtoken = ?", x).First(&user)
	err = SQLrequest.Error
	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println("Database:", dsn, "test connected successfully!")
	return user.Name
}
