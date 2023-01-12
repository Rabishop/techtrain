package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"techtrain/connectdb"
	"techtrain/gacha"
	"techtrain/limitedgacha"
	"techtrain/techdb"
	"techtrain/transfer"
	"time"
)

// UserGetResponse struct
type UserGetResponse struct {
	Status int    `json:"status"`
	Name   string `json:"name"`
}

// UserCreateRequest struct
type UserCreateRequest struct {
	Name string `json:"name"`
}

// UserCreateResponse struct
type UserCreateResponse struct {
	Status int    `json:"status"`
	Xtoken string `json:"xtoken"`
}

// GachaDrawRequest struct
type GachaDrawRequest struct {
	Times      int    `json:"times"`
	Pickup     int    `json:"pickup"`
	PrivateKey string `json:"privatekey"`
}

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

// UserCharacter struct
type UserCharacter struct {
	UserCharacterID string `json:"userCharacterID"`
	CharacterID     string `json:"characterID"`
	Name            string `json:"name"`
	Power           string `json:"power"`
}

// GachaDrawResponse struct
type CharacterListResponse struct {
	Status  int             `json:"status"`
	Results []UserCharacter `json:"results"`
}

// ユーザ情報作成API
func UserCreateHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		res := UserCreateResponse{200, ""}
		jsonbyte, err := json.Marshal(res)
		if err != nil {
			fmt.Println("Marshal failed")
		}
		fmt.Fprintln(w, string(jsonbyte))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}

	var userCreateRequest UserCreateRequest
	json.Unmarshal([]byte(body), &userCreateRequest)

	rand.Seed(time.Now().UnixNano())
	xtoken := strconv.Itoa(rand.Int())
	if userCreateRequest.Name == "" {
		userCreateRequest.Name = "Unnamed"
	}
	responseStatus := connectdb.ConnWriteName(xtoken, userCreateRequest.Name)

	userCreateResponse := UserCreateResponse{responseStatus, xtoken}
	jsonbyte, err := json.Marshal(userCreateResponse)
	if err != nil {
		fmt.Println("Marshal failed")
	}

	//Json return
	fmt.Fprintln(w, string(jsonbyte))
}

// ユーザ情報取得API
func UserGetHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		res := UserGetResponse{200, ""}
		jsonbyte, err := json.Marshal(res)
		if err != nil {
			fmt.Println("Marshal failed")
		}
		fmt.Fprintln(w, string(jsonbyte))
		return
	}

	if r.Header.Values("x-token") == nil {
		fmt.Println("None x-token found")
		return
	}

	// read header
	var userName string
	x_token := strings.Trim(r.Header.Values("x-token")[0], "\"")
	responseStatus := connectdb.ConnReadName(x_token, &userName)
	userGetResponse := UserGetResponse{responseStatus, userName}

	jsonbyte, err := json.Marshal(userGetResponse)
	if err != nil {
		fmt.Println("Marshal failed")
	}
	fmt.Fprintln(w, string(jsonbyte))
}

// ユーザ情報更新API
func UserUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		res := UserGetResponse{200, ""}
		jsonbyte, err := json.Marshal(res)
		if err != nil {
			fmt.Println("Marshal failed")
		}
		fmt.Fprintln(w, string(jsonbyte))
		return
	}

	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	var userGetResponse UserGetResponse
	json.Unmarshal([]byte(body), &userGetResponse)

	// read header
	x_token := strings.Trim(r.Header.Values("x-token")[0], "\"")

	// update name
	userGetResponse.Status = connectdb.ConnUpdateName(x_token, userGetResponse.Name)
	if userGetResponse.Status != 100 {
		userGetResponse.Name = ""
	}

	jsonbyte, err := json.Marshal(userGetResponse)
	if err != nil {
		fmt.Println("Marshal failed")
	}
	fmt.Fprintln(w, string(jsonbyte))
}

// ガチャ実行API
func GachaDrawHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		res := GachaDrawResponse{200, make([]GachaResult, 0)}
		jsonbyte, err := json.Marshal(res)
		if err != nil {
			fmt.Println("Marshal failed")
		}
		fmt.Fprintln(w, string(jsonbyte))
		return
	}

	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	var data GachaDrawRequest
	json.Unmarshal([]byte(body), &data)
	if data.Times == 0 {
		res := GachaDrawResponse{400, make([]GachaResult, 0)}
		jsonbyte, err := json.Marshal(res)
		if err != nil {
			fmt.Println("Marshal failed")
		}
		fmt.Fprintln(w, string(jsonbyte))
		return
	}

	if data.Times > 100000 {
		res := GachaDrawResponse{401, make([]GachaResult, 0)}
		jsonbyte, err := json.Marshal(res)
		if err != nil {
			fmt.Println("Marshal failed")
		}
		fmt.Fprintln(w, string(jsonbyte))
		return
	}

	// store confirmation result
	var confirmation_result int
	court := make(chan int)
	status := transfer.GachaTransfer(data.PrivateKey, uint32(data.Times), court)
	if status != 100 {
		res := GachaDrawResponse{status, make([]GachaResult, 0)}
		jsonbyte, err := json.Marshal(res)
		if err != nil {
			fmt.Println("Marshal failed")
		}
		fmt.Fprintln(w, string(jsonbyte))
		return
	}

	var character_prob_table [gacha.MAX_ID]int
	gacha.ConnReadProb(&character_prob_table, data.Pickup)

	turn := data.Times / 1000
	remain := data.Times % 1000

	// read xtoken
	x_token := strings.Trim(r.Header.Values("x-token")[0], "\"")

	// get GachaDrawResponse
	var res GachaDrawResponse
	res.Status = status

	for turn_ := 1; turn_ <= turn; turn_++ {
		var userinventory []techdb.Userinventory
		gacha.Gacha_t(x_token, character_prob_table, &userinventory, 1000)
		for count := 0; count < 1000; count++ {
			res.Results = append(res.Results, GachaResult{strconv.Itoa(int(userinventory[count].Characterid)),
				userinventory[count].Name, int(userinventory[count].Power)})
		}
		go gacha.Insert_res(x_token, &userinventory, &confirmation_result, 1000, court)
	}

	if remain != 0 {
		var userinventory []techdb.Userinventory
		gacha.Gacha_t(x_token, character_prob_table, &userinventory, remain)
		for count := 0; count < remain; count++ {
			res.Results = append(res.Results, GachaResult{strconv.Itoa(int(userinventory[count].Characterid)),
				userinventory[count].Name, int(userinventory[count].Power)})
		}
		go gacha.Insert_res(x_token, &userinventory, &confirmation_result, remain, court)
	}

	jsonbyte, err := json.Marshal(res)
	if err != nil {
		fmt.Println("Marshal failed")
	}

	// Json return
	fmt.Fprintln(w, string(jsonbyte))
}

// Limited gacha API
func LimitedDrawHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		res := GachaDrawResponse{200, make([]GachaResult, 0)}
		jsonbyte, err := json.Marshal(res)
		if err != nil {
			fmt.Println("Marshal failed")
		}
		fmt.Fprintln(w, string(jsonbyte))
		return
	}

	// read xtoken from header
	x_token := strings.Trim(r.Header.Values("x-token")[0], "\"")

	// read gachaDrawRequest from body
	var gachaDrawRequest GachaDrawRequest
	var gachaDrawResponse limitedgacha.GachaDrawResponse

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	json.Unmarshal([]byte(body), &gachaDrawRequest)

	if gachaDrawRequest.Times == 0 {
		gachaDrawResponse := GachaDrawResponse{400, make([]GachaResult, 0)}
		jsonbyte, err := json.Marshal(gachaDrawResponse)
		if err != nil {
			fmt.Println("Marshal failed")
		}
		fmt.Fprintln(w, string(jsonbyte))
		return
	} else if gachaDrawRequest.Times >= 100000 {
		gachaDrawResponse := GachaDrawResponse{401, make([]GachaResult, 0)}
		jsonbyte, err := json.Marshal(gachaDrawResponse)
		if err != nil {
			fmt.Println("Marshal failed")
		}
		fmt.Fprintln(w, string(jsonbyte))
		return
	}

	// send a transcation to Blockchain and store result to courtChan
	var transferResult int
	courtChan := make(chan int)
	responseStatus := transfer.GachaTransfer(gachaDrawRequest.PrivateKey, uint32(gachaDrawRequest.Times), courtChan)
	if responseStatus != 100 {
		gachaDrawResponse := GachaDrawResponse{responseStatus, make([]GachaResult, 0)}
		jsonbyte, err := json.Marshal(gachaDrawResponse)
		if err != nil {
			fmt.Println("Marshal failed")
		}
		fmt.Fprintln(w, string(jsonbyte))
		return
	}

	// read character probability and number
	characterProbTable := make([]int, limitedgacha.MAX_ID)
	characterNumber := make([]int, limitedgacha.MAX_ID)
	characterProbWithLimit := make([]techdb.Characterprobwithlimit, limitedgacha.MAX_ID)
	limitedgacha.ConnReadProb(&characterProbTable, &characterNumber, &characterProbWithLimit, gachaDrawRequest.Pickup)

	// gacha start and if gacha failed gachaResult = -1
	gachaResult := 1
	userInventory := make([]techdb.Userinventory, gachaDrawRequest.Times)
	limitedgacha.Gacha(x_token, characterProbTable, &characterNumber, &userInventory, gachaDrawRequest.Times, &gachaResult, &gachaDrawResponse)
	go limitedgacha.Insert_res(x_token, &userInventory, &transferResult, &gachaResult, gachaDrawRequest.Times, courtChan)

	// store rollback data
	// numbeRollback := make([]int, limitedgacha.MAX_ID)
	// for i := 0; i < limitedgacha.MAX_ID-1; i++ {
	// 	numberRollback[i] = int(characterProbWithLimit[i].Number) - characterNumber[i+1]
	// }

	// update new character number after gacha
	limitedgacha.Update_number(characterProbWithLimit, characterNumber, gachaResult)

	// if update failed return 406
	if gachaResult == -1 {
		res := GachaDrawResponse{406, make([]GachaResult, 0)}
		jsonbyte, err := json.Marshal(res)
		if err != nil {
			fmt.Println("Marshal failed")
		}
		fmt.Fprintln(w, string(jsonbyte))
		return
	}

	// rollback if transcation is failed
	// go limitedgacha.Character_numberRollback(numberRollback, &transferResult, gachaDrawRequest.Pickup)

	jsonbyte, err := json.Marshal(gachaDrawResponse)
	if err != nil {
		fmt.Println("Marshal failed")
	}

	// Json return
	fmt.Fprintln(w, string(jsonbyte))
}

// ユーザ所持キャラクター一覧取得API
func CharacterListHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		res := CharacterListResponse{200, make([]UserCharacter, 0)}
		jsonbyte, err := json.Marshal(res)
		if err != nil {
			fmt.Println("Marshal failed")
		}
		fmt.Fprintln(w, string(jsonbyte))
		return
	}

	// read header
	x_token := strings.Trim(r.Header.Values("x-token")[0], "\"")

	var userinventory []techdb.Userinventory
	var res CharacterListResponse
	res.Status = gacha.ConnReadInfo(x_token, &userinventory)
	for i := 0; i < len(userinventory); i++ {
		// fmt.Println(characterid, name)
		res.Results = append(res.Results, UserCharacter{strconv.Itoa(int(userinventory[i].Usercharacterid)), strconv.Itoa(int(userinventory[i].Characterid)), userinventory[i].Name,
			strconv.Itoa(int(userinventory[i].Power))})
	}

	// null処理
	if res.Results == nil {
		res.Results = make([]UserCharacter, 0)
	}

	jsonbyte, err := json.Marshal(res)
	if err != nil {
		fmt.Println("Marshal failed")
	}

	//Json return
	fmt.Fprintln(w, string(jsonbyte))
}

func main() {

	// transfer XY
	// transfer.GachaTransfer(uint32(1))

	// // creat new database
	// techdb.ConnCreatTable()

	// // set characterinfo and characterprob
	// techdb.ConnSetInfo()
	// techdb.ConnSetProb()

	// // set limited characterinfo and characterprob
	// limitedgacha.ConnSetInfo()

	// reset the table
	limitedgacha.ConnSetProb()

	// gacha test
	// var character_prob_table [gacha.MAX_ID]int
	// gacha.ConnReadProb(&character_prob_table, 4)
	// var times = 15
	// var userinventory []techdb.Userinventory
	// gacha.Gacha_t("5033260496457778156", character_prob_table, &userinventory, times)
	// for i := 0; i < times; i++ {
	// 	fmt.Println(userinventory[i])
	// }

	// list test
	// var userinventory []techdb.Userinventory
	// gacha.ConnReadInfo("5033260496457778156", &userinventory)
	// fmt.Println(userinventory)

	// // update test
	// connectdb.ConnUpdateName("example001", "Bob")

	// // write test
	// // connectdb.ConnWriteName("example004", "Alex")

	// // read test
	// res := connectdb.ConnReadName("example001")
	// fmt.Println(res)

	// ユーザ情報作成API
	http.HandleFunc("/user/create", UserCreateHandler)

	//ユーザ情報取得API
	http.HandleFunc("/user/get", UserGetHandler)

	//ユーザ情報更新API
	http.HandleFunc("/user/update", UserUpdateHandler)

	//ガチャ実行API
	http.HandleFunc("/gacha/draw", GachaDrawHandler)

	//限定ガチャ実行API
	http.HandleFunc("/gacha/limiteddraw", LimitedDrawHandler)

	// ユーザ所持キャラクター一覧取得API
	http.HandleFunc("/character/list", CharacterListHandler)

	http.ListenAndServe(":8080", nil)
}
