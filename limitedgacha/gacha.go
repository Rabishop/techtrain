package limitedgacha

import (
	"fmt"
	"math/rand"
	"sync"
	"techtrain/techdb"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var lock sync.Mutex
var MAX_ID_ = MAX_ID

// get character_permille table
func ConnReadProb(character_prob_table *[MAX_ID]int, character_number *[MAX_ID]int, characterprobwithlimit *[MAX_ID]techdb.Characterprobwithlimit, listid int) {
	// try to connect db
	dsn := rUsername + ":" + rPassword + "@" + rProtocol + "(" + rAddress + ")" + "/" + rDbname
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		panic(err)
	}

	//select all
	SQLrequest := db.Where("Listid = ?", listid).Find(characterprobwithlimit)
	err = SQLrequest.Error
	if err != nil {
		fmt.Println(err)
	}

	var p_list [MAX_ID]int
	for i := 1; i < MAX_ID; i++ {
		p_list[characterprobwithlimit[i-1].Characterid] = int(characterprobwithlimit[i-1].Prob)
		character_number[i] = int(characterprobwithlimit[i-1].Number)
	}

	for i := 1; i < MAX_ID; i++ {
		character_prob_table[i] = character_prob_table[i-1] + p_list[i]
	}

	return
}

func Gacha_t(x string, character_prob_table [MAX_ID]int, character_number *[MAX_ID]int, userinventory *[]techdb.Userinventory, times int) {
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

				// if character_number[i] == 0 {
				// 	t--
				// 	break
				// }
				// if character_number[i] != 999999999 {
				// 	character_number[i]--
				// }
				break
			}
		}
	}

	*userinventory = res

	sqlDB, _ := db.DB()
	defer sqlDB.Close()
}

func Insert_res(x string, res *[]techdb.Userinventory, confirmation_result *int, times int, court chan int) {

	// check confirmation result is not known now
	if *confirmation_result == 0 {
		for {
			if *confirmation_result == 1 {
				break
			}
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

	// get userid
	var user techdb.User
	SQLrequest := db.Where("Xtoken = ?", x).First(&user)
	err = SQLrequest.Error
	if err != nil {
		fmt.Println(err)
		return
	}

	lock.Lock()
	// get last usercharacterid
	var last techdb.Userinventory
	SQLrequest1 := db.Where("userid = ?", user.Userid).Order("usercharacterid desc").Limit(1).Find(&last)
	err = SQLrequest1.Error
	if err != nil {
		fmt.Println(err)
	}

	// update usercharacterid
	for t := 0; t < times; t++ {
		(*res)[t].Usercharacterid = last.Usercharacterid + uint(t) + 1
	}

	// update characternumber
	SQLrequest2 := db.Create(res)
	err = SQLrequest2.Error
	if err != nil {
		fmt.Println(err)
	}
	lock.Unlock()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()
}

func Update_number(characterprobwithlimit [MAX_ID]techdb.Characterprobwithlimit, character_number *[MAX_ID]int) {
	// try to connect db
	dsn := rUsername + ":" + rPassword + "@" + rProtocol + "(" + rAddress + ")" + "/" + rDbname
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		panic(err)
	}

	fmt.Printf("-----Before------After-----\n")
	for count := 0; count < MAX_ID-1; count++ {
		fmt.Printf("| %10d |", characterprobwithlimit[count].Number)
		fmt.Printf(" %10d |\n", character_number[count+1])
	}
	fmt.Printf("-----Before------After-----\n")

	for count := 0; count < MAX_ID-1; count++ {
		characterprobwithlimit[count].Number = uint(character_number[count+1])
	}

	res := characterprobwithlimit[:MAX_ID-1]

	// insert userinventory
	lock.Lock()
	SQLrequest2 := db.Save(&res)
	err = SQLrequest2.Error
	if err != nil {
		fmt.Println(err)
	}
	lock.Unlock()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()
}

func Character_numberRollback(number_rollback *[MAX_ID]int, confirmation_result *int, Pickup int) {

	// check confirmation result is not known now
	if *confirmation_result == 0 {
		for {
			if *confirmation_result == -11 {
				break
			}
			// wait for confirmation result
		}
	}

	if *confirmation_result == 1 {
		return
	}

	// try to connect db
	dsn := rUsername + ":" + rPassword + "@" + rProtocol + "(" + rAddress + ")" + "/" + rDbname
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		panic(err)
	}

	// insert userinventory
	lock.Lock()

	// get character number now
	var res [MAX_ID - 1]techdb.Characterprobwithlimit
	SQLrequest1 := db.Where("listid = ?", Pickup).Find(&res)
	err = SQLrequest1.Error
	if err != nil {
		fmt.Println(err)
	}

	// update character number
	for i := 0; i < MAX_ID-1; i++ {
		res[i].Number += uint(number_rollback[i])
	}

	SQLrequest := db.Save(&res)
	err = SQLrequest.Error
	if err != nil {
		fmt.Println(err)
	}
	lock.Unlock()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()
}
