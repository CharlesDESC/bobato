package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

var db, _ = sql.Open("mysql", "root:root@tcp(localhost)/Bobato")
var query, _ = db.Query("SELECT * FROM bobato.usr")

type Topic struct {
	ID       int    `json:"ID"`
	Title    string `json:"title"`
	UserID   int    `json:"userID"`
	Content  string `json:"content"`
	Theme    string `json:"theme"`
	Date     string `json:"date"`
	Response Response
}

type Response struct {
	ID      int    `json:"ID"`
	UserID  int    `json:"userID"`
	TopicID int    `json:"TopicID"`
	Date    string `json:"date"`
	Content string `json:"content"`
}

type Usr struct {
	ID           int    `json:"ID"`
	Name         string `json:"name"`
	Password     string `json:"password"`
	Confirmation string `json:"confirmation"`
	Email        string `json:"email"`
	PP           string `json:"PP"`
}

type Session struct {
	Name  string
	Token string
}

// type Register struct {
// 	Name     string `json:"name"`
// 	Password string `json:"password"`
// 	Email    string `json:"email"`
// }

// func getUsrs() []Usr {
// 	var user []Usr
// 	for query.Next() {
// 		var usr Usr
// 		query.Scan(&usr.ID, &usr.Name, &usr.Password, &usr.Email, &usr.PP)
// 		user = append(user, usr)
// 	}
// 	return user
// }

func getTopics() []Topic {
	var query, _ = db.Query("SELECT * FROM bobato.topic")
	var topics []Topic
	for query.Next() {
		var topic Topic
		query.Scan(&topic.ID, &topic.Title, &topic.UserID, &topic.Content, &topic.Theme, &topic.Date)
		topics = append(topics, topic)
	}
	return topics
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
}

func topicsHandler(w http.ResponseWriter, r *http.Request) {
	a, _ := json.Marshal(getTopics())
	fmt.Println(string(a))
	w.Write(a)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
}

func topicHandler(w http.ResponseWriter, r *http.Request) {
	pathID := r.URL.Path
	pathID = path.Base(pathID)
	pathIDint, _ := strconv.Atoi(pathID)
	getTopicsVar := getTopics()
	a, _ := json.Marshal(getTopicsVar[pathIDint-1])
	w.Write(a)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
}

func UsrsHandler(w http.ResponseWriter, r *http.Request) {
	var user []Usr
	for query.Next() {
		var usr Usr
		query.Scan(&usr.ID, &usr.Name, &usr.Password, &usr.Email, &usr.PP)
		user = append(user, usr)
	}
	fmt.Println(user)
	a, _ := json.Marshal(user)
	w.Write(a)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user Usr
	var session Session
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&user)
	fmt.Println(user)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	emailVar := `SELECT ID, NAME, PASSWORD, EMAIL FROM bobato.usr WHERE EMAIL="` + user.Email + `" AND PASSWORD="` + user.Password + `"`
	var getRaw = db.QueryRow(emailVar)
	getRaw.Scan(&user.ID, &user.Name, &user.Password, &user.Email)
	fmt.Println(user)
	if user.ID != 0 {
		session.Token = uuid.New().String()
		session.Name = user.Name
		fmt.Println(session.Token)
		a, _ := json.Marshal(session)
		w.Write(a)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var register Usr
	// var session Session
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&register)
	fmt.Println(register)

	if register.Password == register.Confirmation {
		insert := `INSERT INTO bobato.usr (NAME, PASSWORD, EMAIL) VALUES ("` + register.Name + `","` + register.Password + `","` + register.Email + `");`
		_, err := db.Query(insert)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("password != confirm")
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

}

func createTopicHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var topic Topic
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&topic)
	insert := `INSERT INTO bobato.topic (TOPIC_NAME, CONTENT, THEME, USR_ID) VALUES ("` + topic.Title + `","` + topic.Content + `","` + topic.Theme + `",` + strconv.Itoa(topic.UserID) + `);`
	fmt.Println(insert)
	_, err := db.Query(insert)
	if err != nil {
		fmt.Println(err)
	}

}

func responseTopicHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var response Response
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&response)
	fmt.Println(response)
	insert := `INSERT INTO bobato.response (TOPIC_ID, CONTENT, USR_ID) VALUES (` + strconv.Itoa(response.TopicID) + `,"` + response.Content + `",` + strconv.Itoa(response.UserID) + `);`
	fmt.Println(insert)
	_, err := db.Query(insert)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	http.HandleFunc("/api/", apiHandler)
	http.HandleFunc("/api/login", loginHandler)
	http.HandleFunc("/api/register", registerHandler)
	http.HandleFunc("/api/topics", topicsHandler)
	http.HandleFunc("/api/topics/", topicHandler)
	http.HandleFunc("/api/usr/select", UsrsHandler)
	http.HandleFunc("/api/createTopic", createTopicHandler)
	http.HandleFunc("/api/responseTopic", responseTopicHandler)

	log.Fatal(http.ListenAndServe(":55", nil))
}
