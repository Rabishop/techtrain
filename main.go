package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"techtrain/connectdb"
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

	result := make(map[string]string)
	keys := r.URL.Query()
	for k, v := range keys {
		result[k] = v[0]
		xtoken := v[0]
		name := connectdb.ConnReadName(xtoken)
		res := UserGetResponse{name}

		// fmt.Println(res)
		jsonbyte, err := json.Marshal(res)
		if err != nil {
			fmt.Println("Marshal failed")
		}
		fmt.Fprintln(w, string(jsonbyte))
	}
}

// ユーザ情報更新API
func user_update_handler(w http.ResponseWriter, r *http.Request) {

	result := make(map[string]string)
	keys := r.URL.Query()
	for k, v := range keys {
		result[k] = v[0]
	}
	xtoken := result["x-token"]
	name := result["name"]
	connectdb.ConnUpdateName(xtoken, name)
}

func main() {

	// connectdb.ConnUpdateName("example001", "Alice")
	// connectdb.ConnWriteName("example004", "Alex")
	// res := connectdb.ConnReadName("example001")
	// fmt.Println(res)

	//ユーザ情報作成API
	http.HandleFunc("/user/create", user_create_handler)

	//ユーザ情報取得API
	http.HandleFunc("/user/get", user_get_handler)

	//ユーザ情報更新API
	http.HandleFunc("/user/update", user_update_handler)

	http.ListenAndServe(":8080", nil)
}
