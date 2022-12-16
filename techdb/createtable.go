package techdb

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type User struct {
	Userid uint `gorm:"primaryKey;AUTO_INCREMENT"`
	Xtoken string
	Name   string
}

type Userinventory struct {
	Userid          uint `gorm:"primaryKey"`
	Usercharacterid uint `gorm:"primaryKey"`
	Characterid     uint
	Name            string
	Power           uint
}

type Characterinfo struct {
	Characterid uint `gorm:"primaryKey"`
	Name        string
	Stdpower    uint
}

type Characterprob struct {
	Listid      uint `gorm:"primaryKey"`
	Characterid uint `gorm:"primaryKey"`
	Prob        uint
}

type Characterprobwithlimit struct {
	Listid      uint `gorm:"primaryKey"`
	Characterid uint `gorm:"primaryKey"`
	Prob        uint
	Number      uint
}

// read db
func ConnCreatTable() {
	// try to connect db
	dsn := rUsername + ":" + rPassword + "@" + rProtocol + "(" + rAddress + ")" + "/" + rDbname

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		panic(err)
	}

	// 迁移 schema
	db.AutoMigrate(&User{}, &Userinventory{}, &Characterinfo{}, &Characterprob{}, &Characterprobwithlimit{})

}
