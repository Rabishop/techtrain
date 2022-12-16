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
func user_create_handler(w http.ResponseWriter, r *http.Request) {

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

	var data UserCreateRequest
	json.Unmarshal([]byte(body), &data)

	rand.Seed(time.Now().UnixNano())
	xtoken := strconv.Itoa(rand.Int())
	if data.Name == "" {
		data.Name = "unnamed"
	}
	status := connectdb.ConnWriteName(xtoken, data.Name)

	res := UserCreateResponse{status, xtoken}
	jsonbyte, err := json.Marshal(res)
	if err != nil {
		fmt.Println("Marshal failed")
	}

	//Json return
	fmt.Fprintln(w, string(jsonbyte))
}

// ユーザ情報取得API
func user_get_handler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		res := UserGetResponse{200, ""}
		jsonbyte, err := json.Marshal(res)
		if err != nil {
			fmt.Println("Marshal failed")
		}
		fmt.Fprintln(w, string(jsonbyte))
		return
	}

	// read header
	x_token := strings.Trim(r.Header.Values("x-token")[0], "\"")
	var name string
	status := connectdb.ConnReadName(x_token, &name)
	res := UserGetResponse{status, name}

	// fmt.Println(res)
	jsonbyte, err := json.Marshal(res)
	if err != nil {
		fmt.Println("Marshal failed")
	}
	fmt.Fprintln(w, string(jsonbyte))
}

// ユーザ情報更新API
func user_update_handler(w http.ResponseWriter, r *http.Request) {

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
	// log.Printf("%s", body)
	var data UserGetResponse
	json.Unmarshal([]byte(body), &data)

	// read header
	x_token := strings.Trim(r.Header.Values("x-token")[0], "\"")

	// update name
	status := connectdb.ConnUpdateName(x_token, data.Name)
	if status != 100 {
		data.Name = ""
	}

	res := UserGetResponse{status, data.Name}
	jsonbyte, err := json.Marshal(res)
	if err != nil {
		fmt.Println("Marshal failed")
	}
	fmt.Fprintln(w, string(jsonbyte))
}

// ガチャ実行API
func gacha_draw_handler(w http.ResponseWriter, r *http.Request) {

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

// ガチャ実行API
func limitedgacha_draw_handler(w http.ResponseWriter, r *http.Request) {

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

	// read characterprob table
	var character_prob_table [limitedgacha.MAX_ID]int
	var character_number [limitedgacha.MAX_ID]int
	var characterprobwithlimit [limitedgacha.MAX_ID]techdb.Characterprobwithlimit
	limitedgacha.ConnReadProb(&character_prob_table, &character_number, &characterprobwithlimit, data.Pickup)

	turn := data.Times / 1000
	remain := data.Times % 1000

	// read xtoken
	x_token := strings.Trim(r.Header.Values("x-token")[0], "\"")

	// get GachaDrawResponse
	var res GachaDrawResponse
	res.Status = status

	for turn_ := 1; turn_ <= turn; turn_++ {
		var userinventory []techdb.Userinventory
		limitedgacha.Gacha_t(x_token, character_prob_table, &character_number, &userinventory, 1000)
		for count := 0; count < 1000; count++ {
			res.Results = append(res.Results, GachaResult{strconv.Itoa(int(userinventory[count].Characterid)),
				userinventory[count].Name, int(userinventory[count].Power)})
		}
		go gacha.Insert_res(x_token, &userinventory, &confirmation_result, 1000, court)
	}

	if remain != 0 {
		var userinventory []techdb.Userinventory
		limitedgacha.Gacha_t(x_token, character_prob_table, &character_number, &userinventory, remain)
		for count := 0; count < remain; count++ {
			res.Results = append(res.Results, GachaResult{strconv.Itoa(int(userinventory[count].Characterid)),
				userinventory[count].Name, int(userinventory[count].Power)})
		}
		go gacha.Insert_res(x_token, &userinventory, &confirmation_result, remain, court)
	}

	var number_rollback [limitedgacha.MAX_ID]int
	for i := 0; i < limitedgacha.MAX_ID-1; i++ {
		number_rollback[i] = int(characterprobwithlimit[i].Number) - character_number[i+1]
	}

	go limitedgacha.Update_number(characterprobwithlimit, &character_number, &confirmation_result)
	// go limitedgacha.Character_numberRollback(&number_rollback, &confirmation_result, data.Pickup)

	jsonbyte, err := json.Marshal(res)
	if err != nil {
		fmt.Println("Marshal failed")
	}

	// Json return
	fmt.Fprintln(w, string(jsonbyte))
}

// ユーザ所持キャラクター一覧取得API
func character_list_handler(w http.ResponseWriter, r *http.Request) {

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
	http.HandleFunc("/user/create", user_create_handler)

	//ユーザ情報取得API
	http.HandleFunc("/user/get", user_get_handler)

	//ユーザ情報更新API
	http.HandleFunc("/user/update", user_update_handler)

	//ガチャ実行API
	http.HandleFunc("/gacha/draw", gacha_draw_handler)

	//限定ガチャ実行API
	http.HandleFunc("/gacha/limiteddraw", limitedgacha_draw_handler)

	// ユーザ所持キャラクター一覧取得API
	http.HandleFunc("/character/list", character_list_handler)

	http.ListenAndServe(":8080", nil)
}
