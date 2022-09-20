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
	"time"
)

// UserGetResponse struct
type UserGetResponse struct {
	Name string `json:"name"`
}

// UserCreateRequest struct
type UserCreateRequest struct {
	Name string `json:"name"`
}

// UserCreateResponse struct
type UserCreateResponse struct {
	Xtoken string `json:"xtoken"`
}

// GachaDrawRequest struct
type GachaDrawRequest struct {
	Times int `json:"times"`
}

// GachaResult struct
type GachaResult struct {
	CharacterID string `json:"characterID"`
	Name        string `json:"name"`
}

// GachaDrawResponse struct
type GachaDrawResponse struct {
	Results []GachaResult `json:"results"`
}

// UserCharacter struct
type UserCharacter struct {
	UserCharacterID string `json:"userCharacterID"`
	CharacterID     string `json:"characterID"`
	Name            string `json:"name"`
}

// GachaDrawResponse struct
type CharacterListResponse struct {
	Results []UserCharacter `json:"results"`
}

// ユーザ情報作成API
func user_create_handler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%s", body)
	var data UserCreateRequest
	json.Unmarshal([]byte(body), &data)

	rand.Seed(time.Now().UnixNano())
	xtoken := strconv.Itoa(rand.Int())

	connectdb.ConnWriteName(xtoken, data.Name)

	res := UserCreateResponse{xtoken}
	jsonbyte, err := json.Marshal(res)
	if err != nil {
		fmt.Println("Marshal failed")
	}

	//Json return
	fmt.Fprintln(w, string(jsonbyte))
}

// ユーザ情報取得API
func user_get_handler(w http.ResponseWriter, r *http.Request) {

	// read header
	x_token := strings.Trim(r.Header.Values("x-token")[0], "\"")
	name := connectdb.ConnReadName(x_token)
	res := UserGetResponse{name}

	// fmt.Println(res)
	jsonbyte, err := json.Marshal(res)
	if err != nil {
		fmt.Println("Marshal failed")
	}
	fmt.Fprintln(w, string(jsonbyte))
}

// ユーザ情報更新API
func user_update_handler(w http.ResponseWriter, r *http.Request) {

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
	connectdb.ConnUpdateName(x_token, data.Name)
}

// ガチャ実行API
func gacha_draw_handler(w http.ResponseWriter, r *http.Request) {

	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	// log.Printf("%s", body)
	var data GachaDrawRequest
	json.Unmarshal([]byte(body), &data)

	// read header
	x_token := strings.Trim(r.Header.Values("x-token")[0], "\"")
	// fmt.Println(r.Header.Values("x-token")[0], data.Times)

	var character_permille [1001]int
	gacha.ConnReadProb(&character_permille)

	var res GachaDrawResponse
	for t := 1; t <= data.Times; t++ {
		var characterid string
		var name string
		gacha.Gacha(x_token, character_permille, &characterid, &name)
		// fmt.Println(characterid, name)
		res.Results = append(res.Results, GachaResult{characterid, name})
	}

	fmt.Println(res)

	jsonbyte, err := json.Marshal(res)
	if err != nil {
		fmt.Println("Marshal failed")
	}

	//Json return
	fmt.Fprintln(w, string(jsonbyte))
}

// ユーザ所持キャラクター一覧取得API
func character_list_handler(w http.ResponseWriter, r *http.Request) {

	// read header
	x_token := strings.Trim(r.Header.Values("x-token")[0], "\"")
	// fmt.Println(r.Header.Values("x-token")[0])

	// read character list
	var character_list [gacha.MAX_ID]string
	var user_inventory [gacha.MAX_ID]int
	gacha.ConnReadInfo(&character_list)
	gacha.ConnReadList(x_token, &user_inventory)
	fmt.Println(user_inventory)

	var res CharacterListResponse
	for i := 1; i < gacha.MAX_ID; i++ {
		for j := 1; j <= user_inventory[i]; j++ {
			// fmt.Println(characterid, name)
			res.Results = append(res.Results, UserCharacter{x_token + "_" + strconv.Itoa(i) + "_" + strconv.Itoa(j), strconv.Itoa(i), character_list[i]})
		}
	}

	// fmt.Println(res)

	jsonbyte, err := json.Marshal(res)
	if err != nil {
		fmt.Println("Marshal failed")
	}

	//Json return
	fmt.Fprintln(w, string(jsonbyte))
}

func main() {

	// gacha test
	// var character_permille [1001]int
	// var characterid string
	// var name string
	// gacha.ConnReadProb(&character_permille)
	// for item := 1; item <= 1000; item++ {
	// 	fmt.Printf("i:%d character:%d\n", item, character_permille[item])
	// }
	// gacha.Gacha("example001", character_permille, &characterid, &name)
	// fmt.Println(characterid, name)

	// list test
	// var list [gacha.MAX_ID]int
	// gacha.ConnReadList("example001", &list)
	// fmt.Println(list)
	// var list [gacha.MAX_ID]string
	// gacha.ConnReadInfo(&list)
	// fmt.Println(list)

	// update test
	// connectdb.ConnUpdateName("example001", "Alice")

	// write test
	// connectdb.ConnWriteName("example004", "Alex")

	// read test
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

	// ユーザ所持キャラクター一覧取得API
	http.HandleFunc("/character/list", character_list_handler)

	http.ListenAndServe(":8080", nil)
}
