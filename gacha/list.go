package gacha

import (
	"database/sql"
	"fmt"
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
		xtoken string      `db:"xtoken"`
		id     [MAX_ID]int `db:"id"`
	}
	//select user's inventory
	SQLrequest := "SELECT * FROM userinventory WHERE xtoken = " + "\"" + x + "\""
	// fmt.Println(SQLrequest)
	rows, err := rdb.Query(SQLrequest)
	if err != nil {
		fmt.Println(err)
	}

	// var res string
	for rows.Next() {
		var s info
		err = rows.Scan(&s.xtoken, &s.id[1], &s.id[2], &s.id[3], &s.id[4], &s.id[5], &s.id[6], &s.id[7], &s.id[8], &s.id[9], &s.id[10])
		// fmt.Println(s)
		for item := 1; item < MAX_ID; item++ {
			list[item] = s.id[item]
		}
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
		level       int    `db:"level"`
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
		err = rows.Scan(&s.characterid, &s.name, &s.level)
		// fmt.Println(s)
		character_list[item] = s.name
		item++
	}

	//close rows
	rows.Close()

	// fmt.Println("Database:", dsn, "test connected successfully!")
	return
}
