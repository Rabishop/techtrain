package limitedgacha

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var err error

// MAX_ID >= 1
// データベース内のcharacteridを順番に排列(1 ~ MAX_ID-1)
const MAX_ID = 12

// write db setup
var wdb *sql.DB

const (
	wDriverName = "mysql"

	wUsername = "root"
	wPassword = "kurumi9452"
	wProtocol = "tcp"
	wAddress  = "127.0.0.1:3306"
	wDbname   = "techtrain"
)

// read db setup
var rdb *sql.DB

const (
	rDriverName = "mysql"

	rUsername = "root"
	rPassword = "kurumi9452"
	rProtocol = "tcp"
	rAddress  = "127.0.0.2:3306"
	rDbname   = "techtrain"
)
