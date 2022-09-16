package connectdb

import (
	"database/sql"
	"fmt"
)

func ConnUpdateName(x string, n string) {
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
	SQLrequest := "UPDATE user SET name = \"" + n + "\" Where xtoken=\"" + x + "\""
	fmt.Println(SQLrequest)
	rows, err := wdb.Query(SQLrequest)
	if err != nil {
		fmt.Println(err)
	}

	//print all
	for rows.Next() {
		var s info
		err = rows.Scan(&s.xtoken, &s.name)
		// fmt.Println(s)
	}

	//close rows
	rows.Close()

	fmt.Println("Database:", dsn, "test connected successfully!")
	return
}
