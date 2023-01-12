package limitedgacha

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"techtrain/techdb"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// GachaResult struct
type GachaResult struct {
	CharacterID string `json:"characterID"`
	Name        string `json:"name"`
	Power       int    `json:"power"`
}

// GachaDrawResponse struct
type GachaDrawResponse struct {
	Status  int           `json:"status"`
	Results []GachaResult `json:"results"`
}

var lock sync.Mutex
var MAX_ID_ = MAX_ID

// get character_permille table
func ConnReadProb(character_prob_table *[]int, character_number *[]int, characterprobwithlimit *[]techdb.Characterprobwithlimit, listid int) {
	// try to connect db
	dsn := Username + ":" + Password + "@" + Protocol + "(" + Address + ")" + "/" + Dbname
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		panic(err)
	}

	//select all
	lock.Lock()
	SQLrequest := db.Where("Listid = ?", listid).Find(characterprobwithlimit)
	err = SQLrequest.Error
	if err != nil {
		fmt.Println(err)
	}
	lock.Unlock()

	var p_list [MAX_ID]int
	for i := 1; i < MAX_ID; i++ {
		p_list[(*characterprobwithlimit)[i-1].Characterid] = int((*characterprobwithlimit)[i-1].Prob)
		(*character_number)[i] = int((*characterprobwithlimit)[i-1].Number)
	}

	for i := 1; i < MAX_ID; i++ {
		(*character_prob_table)[i] = (*character_prob_table)[i-1] + p_list[i]
	}

	return
}

func Gacha(xtoken string, characterProbTable []int, characterNumber *[]int, userInventory *[]techdb.Userinventory, times int, gachaResult *int, gachaDrawResponse *GachaDrawResponse) {
	Gacha_t(xtoken, characterProbTable, characterNumber, userInventory, times, gachaResult)
	for count := 0; count < times; count++ {
		gachaDrawResponse.Results = append(gachaDrawResponse.Results, GachaResult{strconv.Itoa(int((*userInventory)[count].Characterid)),
			(*userInventory)[count].Name, int((*userInventory)[count].Power)})
	}
}

func Gacha_t(x string, character_prob_table []int, characterNumber *[]int, userinventory *[]techdb.Userinventory, times int, gachalimit_result *int) {
	// try to connect db
	dsn := Username + ":" + Password + "@" + Protocol + "(" + Address + ")" + "/" + Dbname

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

				if (*characterNumber)[i] == 0 {
					*gachalimit_result = -1
					t--
				}
				if (*characterNumber)[i] != 999999999 {
					(*characterNumber)[i]--
				}
				break
			}
		}
	}

	*userinventory = res

	sqlDB, _ := db.DB()
	defer sqlDB.Close()
}

func Insert_res(x string, res *[]techdb.Userinventory, confirmation_result *int, gachalimit_result *int, times int, court chan int) {

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

	if *gachalimit_result == -1 {
		return
	}

	// try to connect db
	dsn := Username + ":" + Password + "@" + Protocol + "(" + Address + ")" + "/" + Dbname
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

func Update_number(characterprobwithlimit []techdb.Characterprobwithlimit, character_number []int, gachalimit_result int) {
	// try to connect db
	dsn := Username + ":" + Password + "@" + Protocol + "(" + Address + ")" + "/" + Dbname
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		panic(err)
	}

	if gachalimit_result == -1 {
		fmt.Println("Gacha failed: Limited character run out")
		fmt.Printf("-----Before------After-----\n")
		for count := 0; count < MAX_ID-1; count++ {
			fmt.Printf("| %10d |", characterprobwithlimit[count].Number)
			fmt.Printf(" %10d |\n", characterprobwithlimit[count].Number)
		}
		fmt.Printf("-----Before------After-----\n")
		return
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

func Character_numberRollback(number_rollback []int, confirmation_result *int, Pickup int) {

	// check confirmation result is not known now
	if *confirmation_result == 0 {
		for {
			if *confirmation_result == -1 {
				break
			}
			if *confirmation_result == 1 {
				return
			}
			// wait for confirmation result
		}
	}

	// try to connect db
	dsn := Username + ":" + Password + "@" + Protocol + "(" + Address + ")" + "/" + Dbname
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
