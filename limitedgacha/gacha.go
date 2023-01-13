package limitedgacha

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
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

func Gacha(xtoken string, times int, gachaDrawResponse *GachaDrawResponse, courtChan chan int) error {

	// connect database
	dsn := Username + ":" + Password + "@" + Protocol + "(" + Address + ")" + "/" + Dbname
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		panic(err)
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// start transcation
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// initialize gachaResult,transferResult and userInventory
	gachaResult := 1
	transferResult := 0
	userInventory := make([]techdb.Userinventory, times)

	// read character probability and number
	characterProbTable := make([]int, MAX_ID)
	characterNumber := make([]int, MAX_ID)
	numberRollback := make([]int, MAX_ID)
	characterProbWithLimit := make([]techdb.Characterprobwithlimit, MAX_ID)
	// for i := 0; i < limitedgacha.MAX_ID-1; i++ {
	// 	numberRollback[i] = int(characterProbWithLimit[i].Number) - characterNumber[i+1]
	// }

	if err := ReadProbAndNumber(tx, &characterProbTable, &characterNumber, &characterProbWithLimit, 1); err != nil {
		tx.Rollback()
		return err
	}

	if err := Draw(tx, xtoken, characterProbTable, &characterNumber, &userInventory, gachaDrawResponse, times, &gachaResult); err != nil {
		tx.Rollback()
		return err
	}

	if err := UpdateNumber(tx, characterProbWithLimit, characterNumber, &numberRollback, gachaResult); err != nil {
		tx.Rollback()
		return err
	}

	go InsertResult(xtoken, &userInventory, &transferResult, &gachaResult, times, courtChan)

	go characterNumberRollback(numberRollback, &transferResult, &gachaResult)

	return tx.Commit().Error
}

func ReadProbAndNumber(tx *gorm.DB, characterProbTable *[]int, characterNumber *[]int, characterProbWithLimit *[]techdb.Characterprobwithlimit, listid int) error {
	// try to connect db

	if err := tx.Where("Listid = ?", listid).Find(characterProbWithLimit).Error; err != nil {
		return err
	}

	probList := make([]int, MAX_ID)
	for i := 1; i < MAX_ID; i++ {
		probList[(*characterProbWithLimit)[i-1].Characterid] = int((*characterProbWithLimit)[i-1].Prob)
		(*characterNumber)[i] = int((*characterProbWithLimit)[i-1].Number)
	}

	for i := 1; i < MAX_ID; i++ {
		(*characterProbTable)[i] = (*characterProbTable)[i-1] + probList[i]
	}

	return nil
}

func Draw(tx *gorm.DB, x string, characterProbTable []int, characterNumber *[]int, userInventory *[]techdb.Userinventory, gachaDrawResponse *GachaDrawResponse, times int, gachaResult *int) error {

	// read userId
	var user techdb.User
	if err := tx.Where("Xtoken = ?", x).First(&user).Error; err != nil {
		return err
	}

	// read last userCharacterId
	var lastCharacter techdb.Userinventory
	if err := tx.Where("userid = ?", user.Userid).Order("usercharacterid desc").Limit(1).Find(&lastCharacter).Error; err != nil {
		return err
	}

	// read character name and power
	var characterinfo [MAX_ID]techdb.Characterinfo
	if err := tx.Find(&characterinfo).Error; err != nil {
		return err
	}

	// draw card
	rand.Seed(time.Now().UnixNano())
	for t := 0; t < times; t++ {
		nonce := rand.Intn(100000) + 1
		for i := 1; i < MAX_ID; i++ {
			if nonce <= characterProbTable[i] {
				(*userInventory)[t].Userid = user.Userid
				(*userInventory)[t].Characterid = uint(i)
				(*userInventory)[t].Name = characterinfo[i-1].Name
				(*userInventory)[t].Power = characterinfo[i-1].Stdpower + uint(nonce)/250 - 200

				if (*characterNumber)[i] == 0 {
					*gachaResult = -1
				}
				if (*characterNumber)[i] != 999999999 {
					(*characterNumber)[i]--
				}
				break
			}
		}

		gachaDrawResponse.Results = append(gachaDrawResponse.Results, GachaResult{strconv.Itoa(int((*userInventory)[t].Characterid)),
			(*userInventory)[t].Name, int((*userInventory)[t].Power)})
	}

	return nil
}

func UpdateNumber(tx *gorm.DB, characterProbWithLimit []techdb.Characterprobwithlimit, characterNumber []int, numberRollback *[]int, gachaResult int) error {

	fmt.Printf("-----Before------After-----\n")
	for count := 0; count < MAX_ID-1; count++ {
		fmt.Printf("| %10d |", characterProbWithLimit[count].Number)
		fmt.Printf(" %10d |\n", characterNumber[count+1])
	}
	fmt.Printf("-----Before------After-----\n")
	if gachaResult == -1 {
		return errors.New("Gacha failed: Limited character run out...")
	}

	for count := 0; count < MAX_ID-1; count++ {
		characterProbWithLimit[count].Number = uint(characterNumber[count+1])
		(*numberRollback)[count] = int(characterProbWithLimit[count].Number) - characterNumber[count+1]
	}

	updatedCharacterProbWithLimit := characterProbWithLimit[:MAX_ID-1]

	// insert userinventory
	if err := tx.Save(&updatedCharacterProbWithLimit).Error; err != nil {
		return err
	}

	return nil
}

func InsertResult(x string, userInventory *[]techdb.Userinventory, transferResult *int, gachaResult *int, times int, courtChan chan int) {

	// connect database
	dsn := Username + ":" + Password + "@" + Protocol + "(" + Address + ")" + "/" + Dbname
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		panic(err)
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// start transcation
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// check transfer result
	if *transferResult == 0 {
		for {
			if *transferResult != 0 {
				break
			}
			// wait for confirmation result
			result, ok := <-courtChan
			if ok {
				*transferResult = result
				close(courtChan)
				break
			}
		}
	}

	if *transferResult == -1 || *gachaResult == -1 {
		return
	}

	// read userid
	var user techdb.User
	if err := tx.Where("Xtoken = ?", x).First(&user).Error; err != nil {
		fmt.Println(err)
		return
	}

	// get last usercharacterid
	var lastCharacter techdb.Userinventory
	if err := tx.Where("userid = ?", user.Userid).Order("usercharacterid desc").Limit(1).Find(&lastCharacter).Error; err != nil {
		fmt.Println(err)
		return
	}

	// update usercharacterid
	for t := 0; t < times; t++ {
		(*userInventory)[t].Usercharacterid = lastCharacter.Usercharacterid + uint(t) + 1
	}

	// update characternumber
	for start := 0; start < len(*userInventory); start += 10000 {
		end := start + 10000
		if end > len(*userInventory) {
			end = len(*userInventory)
		}
		batch := (*userInventory)[start:end]
		if err := tx.Create(&batch).Error; err != nil {
			tx.Rollback()
			fmt.Println(err)
			return
		}
	}

	tx.Commit()
}

func characterNumberRollback(numberRollback []int, transferResult *int, gachaResult *int) {

	if *transferResult == 0 {
		for {
			if *transferResult != 0 {
				break
			}
		}
	}

	if *transferResult == -1 && *gachaResult == 1 {
		fmt.Println("Recovering the character number...")
	} else {
		return
	}

	// connect database
	dsn := Username + ":" + Password + "@" + Protocol + "(" + Address + ")" + "/" + Dbname
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		panic(err)
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// start transcation
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// get character number now
	var characterProbWithLimit [MAX_ID - 1]techdb.Characterprobwithlimit
	if err := tx.Where("listid = ?", 1).Find(&characterProbWithLimit).Error; err != nil {
		fmt.Println(err)
		return
	}

	// update character number
	for i := 0; i < MAX_ID-1; i++ {
		characterProbWithLimit[i].Number += uint(numberRollback[i])
	}

	if err := tx.Save(&characterProbWithLimit).Error; err != nil {
		tx.Rollback()
		fmt.Println(err)
		return
	}

	tx.Commit()
}
