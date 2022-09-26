package gacha

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// get character_permille table
func ConnReadProb(character_prob_table *[MAX_ID]int, pickup int) {
	// try to connect db
	dsn := rUsername + ":" + rPassword + "@" + rProtocol + "(" + rAddress + ")" + "/" + rDbname + "?multiStatements=true"
	rdb, err = sql.Open(rDriverName, dsn)
	if err != nil {
		panic(err)
	}

	//table struct
	type prob struct {
		characterid int `db:"characterid"`
		prob        int `db:"prob"`
	}

	//check pickup and access SQL
	SQLrequest := "SELECT * FROM characterprob"
	if pickup != 0 {
		SQLrequest += strconv.Itoa(pickup)
	}

	rows, err := rdb.Query(SQLrequest)
	if err != nil {
		fmt.Println(err)
	}

	var p_list [MAX_ID]int

	//get characterid and prob list
	for rows.Next() {
		var s prob
		err = rows.Scan(&s.characterid, &s.prob)
		// fmt.Println(s.characterid, s.prob)

		p_list[s.characterid] = s.prob
		// fmt.Println(c_list[characterid_count], p_list[characterid_count])
	}

	character_prob_table[1] = p_list[1]
	for i := 2; i < MAX_ID; i++ {
		character_prob_table[i] = character_prob_table[i-1] + p_list[i]
	}

	//close rows
	rows.Close()
	// fmt.Println("Database:", dsn, "test connected successfully!")

	return
}

func Gacha_t(x string, character_prob_table [MAX_ID]int, characterid *[1001]string, name *[1001]string, times int) {
	// draw by rand num
	var id_number [MAX_ID]int
	rand.Seed(time.Now().UnixNano())
	for t := 1; t <= times; t++ {
		nonce := rand.Intn(100000) + 1
		for i := 1; i < MAX_ID; i++ {
			if nonce <= character_prob_table[i] {
				characterid[t] = strconv.Itoa(i)
				id_number[i]++
				break
			}
		}
	}

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

	// get character name from mysql
	SQLrequest1 := "SELECT * FROM characterinfo WHERE characterid = " + characterid[1]

	for t := 2; t <= times; t++ {
		SQLrequest1 += "\nUNION ALL\nSELECT * FROM characterinfo WHERE characterid = " + characterid[t]
	}

	// fmt.Println(SQLrequest1)
	rows, err1 := rdb.Query(SQLrequest1)
	if err1 != nil {
		fmt.Println(err1)
	}

	// get charactername
	name_count := 1
	for rows.Next() {
		var s info
		err = rows.Scan(&s.characterid, &s.name)
		name[name_count] = s.name
		name_count++
		// fmt.Println(s.characterid, s.name)
	}

	rows.Close()

	// update userinventory
	for i := 1; i < MAX_ID; i++ {
		if id_number[i] == 0 {
			continue
		}
		SQLrequest2 := "UPDATE userinventory SET number = number + " + strconv.Itoa(id_number[i]) +
			" Where xtoken=\"" + x + "\" and characterid = " + strconv.Itoa(i)

		// fmt.Println(SQLrequest2)
		rows, err2 := rdb.Query(SQLrequest2)
		if err2 != nil {
			fmt.Println(err2)
		}

		rows.Close()
	}

	// fmt.Println("Database:", dsn, "test connected successfully!")
}
