package limitedgacha

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var err error

// MAX_ID >= 1
// データベース内のcharacteridを順番に排列(1 ~ MAX_ID-1)
const MAX_ID = 12

var db *sql.DB

const (
	DriverName = "mysql"

	Username = "root"
	Password = "kurumi9452"
	Protocol = "tcp"
	Address  = "127.0.0.2:3306"
	Dbname   = "techtrain"
)
