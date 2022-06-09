package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"
)

type Topic struct {
	Title    string `json:"title"`
	Date     string `json:"date"`
	Author   string `json:"Author"`
	Content  string `json:"content"`
	Response int    `json:"response"`
	//Theme  string `json:"theme"`
}

func getTopics() []Topic {
	return []Topic{
		{
			Title:    "Oué ya un pelo il a tapai mon bato",
			Date:     "32/01/1552",
			Author:   "Pouetos",
			Content:  "gloubiboulga",
			Response: 3,
		},
		{
			Title:    "xccvsvzjdjchv s j",
			Date:     "32/01/1552",
			Author:   "Pouetos",
			Content:  "gloubibqsdqsdoulga",
			Response: 3,
		},
		{
			Title:    "nzefinjfez",
			Date:     "32/01/1552",
			Author:   "Pouetos",
			Content:  "gloubibouaéoiépé1121212lga",
			Response: 3,
		},
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
}

func topicsHandler(w http.ResponseWriter, r *http.Request) {

	a, _ := json.Marshal(getTopics())
	w.Write(a)
}

func topicHandler(w http.ResponseWriter, r *http.Request) {
	pathID := r.URL.Path
	pathID = path.Base(pathID)
	pathIDint, _ := strconv.Atoi(pathID)
	getTopicsVar := getTopics()
	fmt.Println(pathIDint)
	a, _ := json.Marshal(getTopicsVar[pathIDint-1])
	w.Write(a)
}

func main() {
	http.HandleFunc("/api/", apiHandler)
	// http.HandleFunc("/api/login", loginHandler)
	// http.HandleFunc("/api/register", registerHandler)
	http.HandleFunc("/api/topics", topicsHandler)
	http.HandleFunc("/api/topics/", topicHandler)

	log.Fatal(http.ListenAndServe(":55", nil))
}
