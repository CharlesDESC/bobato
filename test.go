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
	Name     string `json:"name"`
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
	UserID int
	Token  string
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
		topic.Response = getResponse(topic.ID)
		// fmt.Println(topic.ID)
		var queryName = db.QueryRow(`SELECT NAME FROM bobato.usr WHERE ID=` + strconv.Itoa(topic.ID))
		queryName.Scan(&topic.Name)
		topics = append(topics, topic)
	}

	return topics
}
func getResponse(id int) Response {
	var queryResponse = db.QueryRow(`SELECT ID, USR_ID, TOPIC_ID, CONTENT FROM bobato.response WHERE TOPIC_ID=` + strconv.Itoa(id))
	// fmt.Println(queryResponse)
	var response Response
	queryResponse.Scan(&response.ID, &response.UserID, &response.TopicID, &response.Content)
	// fmt.Println(response)
	return response
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
}

func topicsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	a, _ := json.Marshal(getTopics())
	// fmt.Println(string(a))
	w.Write(a)

}

func topicHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	pathID := r.URL.Path
	pathID = path.Base(pathID)
	pathIDint, _ := strconv.Atoi(pathID)
	getTopicsVar := getTopics()
	a, _ := json.Marshal(getTopicsVar[pathIDint-1])
	w.Write(a)
}

func UsrsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var user []Usr
	for query.Next() {
		var usr Usr
		query.Scan(&usr.ID, &usr.Name, &usr.Password, &usr.Email, &usr.PP)
		user = append(user, usr)
	}
	// fmt.Println(user)
	a, _ := json.Marshal(user)
	w.Write(a)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var user Usr
	var session Session
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&user)
	// fmt.Println(user)
	emailVar := `SELECT ID, NAME, PASSWORD, EMAIL FROM bobato.usr WHERE EMAIL="` + user.Email + `" AND PASSWORD="` + user.Password + `"`
	var getRaw = db.QueryRow(emailVar)
	getRaw.Scan(&user.ID, &user.Name, &user.Password, &user.Email)
	// fmt.Println(user)
	if user.ID != 0 {
		session.Token = uuid.New().String()
		session.UserID = user.ID
		// fmt.Println(session.Token)
		insert := `INSERT INTO bobato.session (USR_ID, TOKEN) VALUES (` + strconv.Itoa(session.UserID) + `,"` + session.Token + `");`
		// fmt.Println(insert)
		_, err := db.Query(insert)
		a, _ := json.Marshal(session)
		w.Write(a)
		if err != nil {
			// fmt.Println(err)
		}

	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var register Usr
	// var session Session
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&register)
	// fmt.Println(register)

	if register.Password == register.Confirmation {
		insert := `INSERT INTO bobato.usr (NAME, PASSWORD, EMAIL) VALUES ("` + register.Name + `","` + register.Password + `","` + register.Email + `");`
		_, err := db.Query(insert)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// fmt.Println("password != confirm")
	}

}

func createTopicHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var topic Topic
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&topic)
	fmt.Println(topic)
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
	// fmt.Println(response)
	insert := `INSERT INTO bobato.response (TOPIC_ID, CONTENT) VALUES (` + strconv.Itoa(response.TopicID) + `,"` + response.Content + `",` + `);`
	// fmt.Println(insert)
	_, err := db.Query(insert)
	if err != nil {
		// fmt.Println(err)
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
