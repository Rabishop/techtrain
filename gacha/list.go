package gacha

import (
	"database/sql"
	"fmt"
	"strconv"
)

func ConnReadList(x string, list *[MAX_ID]int) {
	// try to connect db
	dsn := rUsername + ":" + rPassword + "@" + rProtocol + "(" + rAddress + ")" + "/" + rDbname
	rdb, err = sql.Open(rDriverName, dsn)
	if err != nil {
		panic(err)
	}

	//table struct
	type info struct {
		xtoken      string `db:"xtoken"`
		characterid int    `db:"characterid"`
		number      int    `db:"number"`
	}
	//select user's inventory
	SQLrequest := "SELECT * FROM userinventory WHERE xtoken = " + "\"" + x + "\" and characterid = 1"
	for i := 2; i < MAX_ID; i++ {
		SQLrequest += "\nUNION\nSELECT * FROM userinventory WHERE xtoken = " + "\"" + x + "\" and characterid = " + strconv.Itoa(i)
	}

	// fmt.Println(SQLrequest)
	rows, err := rdb.Query(SQLrequest)
	if err != nil {
		fmt.Println(err)
	}

	// var res string
	for rows.Next() {
		var s info
		err = rows.Scan(&s.xtoken, &s.characterid, &s.number)
		// fmt.Println(s)
		list[s.characterid] = s.number
	}

	//close rows
	rows.Close()

	// fmt.Println("Database:", dsn, "test connected successfully!")
	return
}

func ConnReadInfo(character_list *[MAX_ID]string) {
	// try to connect db
	dsn := rUsername + ":" + rPassword + "@" + rProtocol + "(" + rAddress + ")" + "/" + rDbname
	rdb, err = sql.Open(rDriverName, dsn)
	if err != nil {
		panic(err)
	}

	//table struct
	type info struct {
		characterid int    `db:"characterid"`
		name        string `db:"name"`
	}
	//select user's inventory
	SQLrequest := "SELECT * FROM characterinfo"
	// fmt.Println(SQLrequest)
	rows, err := rdb.Query(SQLrequest)
	if err != nil {
		fmt.Println(err)
	}

	// get character list
	item := 1
	for rows.Next() {
		var s info
		err = rows.Scan(&s.characterid, &s.name)
		// fmt.Println(s)
		character_list[item] = s.name
		item++
	}

	//close rows
	rows.Close()

	// fmt.Println("Database:", dsn, "test connected successfully!")
	return
}
