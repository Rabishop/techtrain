package gacha

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// get character_permille table
func ConnReadProb(character_permille *[1001]int) {
	// try to connect db
	dsn := rUsername + ":" + rPassword + "@" + rProtocol + "(" + rAddress + ")" + "/" + rDbname
	rdb, err = sql.Open(rDriverName, dsn)
	if err != nil {
		panic(err)
	}

	//table struct
	type prob struct {
		characterid int `db:"characterid"`
		prob        int `db:"prob"`
	}
	//select all
	rows, err := rdb.Query("SELECT * FROM characterprob")
	if err != nil {
		fmt.Println(err)
	}

	var c_list [MAX_ID]int
	var p_list [MAX_ID]int

	//get characterid and prob list
	characterid_count := 1
	for rows.Next() {
		var s prob
		err = rows.Scan(&s.characterid, &s.prob)
		// fmt.Println(s.characterid, s.prob)

		c_list[characterid_count] = s.characterid
		p_list[characterid_count] = s.prob
		// fmt.Println(c_list[characterid_count], p_list[characterid_count])
		characterid_count++
	}

	//character per-mille table
	var row_now int = 1
	var prob_now int = p_list[1]
	for item := 1; item <= 1000; item++ {
		if item <= prob_now {
			character_permille[item] = c_list[row_now]
		} else {
			row_now++
			prob_now += p_list[row_now]
			character_permille[item] = c_list[row_now]
		}
		// fmt.Printf("i:%d character:%d\n", item, character_permille[item])
	}

	//close rows
	rows.Close()
	// fmt.Println("Database:", dsn, "test connected successfully!")

	return
}

func Gacha(x string, character_permille [1001]int, characterid *string, name *string) {
	// draw by rand num
	rand.Seed(time.Now().UnixNano())
	permille_id := rand.Intn(1000) + 1
	*characterid = strconv.Itoa(character_permille[permille_id])

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

	// get character name from mysql
	SQLrequest1 := "SELECT * FROM techtrain.characterinfo WHERE characterid = " + *characterid
	// fmt.Println(SQLrequest1)
	rows, err1 := rdb.Query(SQLrequest1)
	if err1 != nil {
		fmt.Println(err1)
	}

	// get characterid and prob list
	for rows.Next() {
		var s info
		err = rows.Scan(&s.characterid, &s.name, &s.level)
		*name = s.name
		// fmt.Println(s.characterid, s.name)
	}

	// update userinventory
	SQLrequest2 := "UPDATE userinventory SET id" + *characterid + " = id" + *characterid + " + 1 Where xtoken=\"" + x + "\""
	// fmt.Println(SQLrequest2)
	rows, err2 := rdb.Query(SQLrequest2)
	if err2 != nil {
		fmt.Println(err2)
	}

	//close rows
	rows.Close()
	// fmt.Println("Database:", dsn, "test connected successfully!")
}
