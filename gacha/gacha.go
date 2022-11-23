package gacha

import (
	"fmt"
	"math/rand"
	"techtrain/techdb"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// get character_permille table
func ConnReadProb(character_prob_table *[11]int, listid int) {
	// try to connect db
	dsn := rUsername + ":" + rPassword + "@" + rProtocol + "(" + rAddress + ")" + "/" + rDbname
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		panic(err)
	}

	//select all
	var characterprob [MAX_ID]techdb.Characterprob
	SQLrequest := db.Where("Listid = ?", listid).Find(&characterprob)
	err = SQLrequest.Error
	if err != nil {
		fmt.Println(err)
	}

	var p_list [MAX_ID]int
	for i := 1; i < MAX_ID; i++ {
		p_list[characterprob[i-1].Characterid] = int(characterprob[i-1].Prob)
	}

	for i := 1; i < MAX_ID; i++ {
		character_prob_table[i] = character_prob_table[i-1] + p_list[i]
	}

	return
}

func Gacha_t(x string, character_prob_table [MAX_ID]int, userinventory *[]techdb.Userinventory, times int) {
	// try to connect db
	dsn := rUsername + ":" + rPassword + "@" + rProtocol + "(" + rAddress + ")" + "/" + rDbname
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		panic(err)
	}

	// get userid
	var user techdb.User
	SQLrequest := db.Where("Xtoken = ?", x).First(&user)
	err = SQLrequest.Error
	if err != nil {
		fmt.Println(err)
		return
	}

	// get last usercharacterid
	var last techdb.Userinventory
	SQLrequest = db.Where("userid = ?", user.Userid).Order("usercharacterid desc").Limit(1).Find(&last)
	err = SQLrequest.Error
	if err != nil {
		fmt.Println(err)
	}

	// get character name and power
	var characterinfo [MAX_ID]techdb.Characterinfo
	SQLrequest = db.Find(&characterinfo)
	err = SQLrequest.Error
	if err != nil {
		fmt.Println(err)
	}

	// draw by rand num
	res := make([]techdb.Userinventory, times)
	rand.Seed(time.Now().UnixNano())
	for t := 0; t < times; t++ {
		nonce := rand.Intn(100000) + 1
		nonce2 := rand.Intn(100000) + 1
		for i := 1; i < MAX_ID; i++ {
			if nonce <= character_prob_table[i] {
				res[t].Userid = user.Userid
				res[t].Usercharacterid = last.Usercharacterid + uint(t) + 1
				res[t].Characterid = uint(i)
				res[t].Name = characterinfo[i-1].Name
				res[t].Power = characterinfo[i-1].Stdpower + uint(nonce2)/250 - 200
				break
			}
		}
	}

	// insert userinventory
	// SQLrequest = db.Create(res)
	// err = SQLrequest.Error
	// if err != nil {
	// 	fmt.Println(err)
	// }

	*userinventory = res

	sqlDB, _ := db.DB()
	defer sqlDB.Close()
}

func Insert_res(res *[]techdb.Userinventory, confirmation_result *int, court chan int) {

	// check confirmation result is not known now
	if *confirmation_result == 0 {
		for {
			// wait for confirmation result
			result, ok := <-court
			if ok {
				*confirmation_result = result
				// fmt.Println(confirmation_result)
				close(court)
				break
			}
		}
	}

	if *confirmation_result == -1 {
		return
	}

	// try to connect db
	dsn := rUsername + ":" + rPassword + "@" + rProtocol + "(" + rAddress + ")" + "/" + rDbname
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		panic(err)
	}

	// insert userinventory
	SQLrequest := db.Create(res)
	err = SQLrequest.Error
	if err != nil {
		fmt.Println(err)
	}

	sqlDB, _ := db.DB()
	defer sqlDB.Close()
}
