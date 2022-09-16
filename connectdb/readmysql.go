package connectdb

import (
	"database/sql"
	"fmt"
)

// read db
func ConnReadMySQL() {
	// try to connect db
	dsn := rUsername + ":" + rPassword + "@" + rProtocol + "(" + rAddress + ")" + "/" + rDbname
	rdb, err = sql.Open(rDriverName, dsn)
	if err != nil {
		panic(err)
	}

	//table struct
	type info struct {
		xtoken string `db:"xtoken"`
		name   string `db:"name"`
	}
	//select all
	rows, err := rdb.Query("SELECT * FROM user")
	if err != nil {
		fmt.Println(err)
		return
	}

	//print all
	for rows.Next() {
		var s info
		err = rows.Scan(&s.xtoken, &s.name)
		fmt.Println(s)
	}

	//close rows
	rows.Close()

	fmt.Println("Database:", dsn, "test connected successfully!")
}

func ConnReadName(x string) string {
	// try to connect db
	dsn := rUsername + ":" + rPassword + "@" + rProtocol + "(" + rAddress + ")" + "/" + rDbname
	rdb, err = sql.Open(rDriverName, dsn)
	if err != nil {
		panic(err)
	}

	//table struct
	type info struct {
		xtoken string `db:"xtoken"`
		name   string `db:"name"`
	}
	//select all
	SQLrequest := "SELECT * FROM techtrain.user WHERE xtoken = " + "\"" + x + "\""
	// fmt.Println(SQLrequest)
	rows, err := rdb.Query(SQLrequest)
	if err != nil {
		fmt.Println(err)
	}

	//print all
	var res string
	for rows.Next() {
		var s info
		err = rows.Scan(&s.xtoken, &s.name)
		// fmt.Println(s)
		res = s.name
	}

	//close rows
	rows.Close()

	fmt.Println("Database:", dsn, "test connected successfully!")
	return res
}
