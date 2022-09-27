package gormuser

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

// read db
func ConnReadMySQL() {
	// try to connect db
	dsn := rUsername + ":" + rPassword + "@" + rProtocol + "(" + rAddress + ")" + "/" + rDbname

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	//table struct
	type info struct {
		xtoken string `db:"xtoken"`
		name   string `db:"name"`
	}

	fmt.Println("Database:", dsn, "test connected successfully!")
}
