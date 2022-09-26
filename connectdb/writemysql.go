package connectdb

import (
	"database/sql"
	"fmt"
	"strconv"
	"techtrain/gacha"
)

// write db
func ConnWriteName(x string, n string) {
	// try to connect db
	dsn := wUsername + ":" + wPassword + "@" + wProtocol + "(" + wAddress + ")" + "/" + wDbname
	wdb, err = sql.Open(rDriverName, dsn)
	if err != nil {
		panic(err)
	}

	//table struct
	type info struct {
		xtoken string `db:"xtoken"`
		name   string `db:"name"`
	}
	//insert xtoken and name
	SQLrequest1 := "INSERT INTO user VALUES (" + "\"" + x + "\",\"" + n + "\")"
	// fmt.Println(SQLrequest1)
	rows, err := wdb.Query(SQLrequest1)
	if err != nil {
		fmt.Println(err)
	}

	//print all
	for rows.Next() {
		var s info
		err = rows.Scan(&s.xtoken, &s.name)
		// fmt.Println(s)
	}

	SQLrequest2 := "INSERT INTO userinventory VALUES (" + "\"" + x + "\",1,0)"
	for item := 2; item < gacha.MAX_ID; item++ {
		SQLrequest2 += ",(" + "\"" + x + "\"," + strconv.Itoa(item) + ",0)"
	}
	// fmt.Println(SQLrequest2)
	rows, err2 := wdb.Query(SQLrequest2)
	if err2 != nil {
		fmt.Println(err2)
	}

	//close rows
	rows.Close()

	// fmt.Println("Database:", dsn, "test connected successfully!")
	return
}
