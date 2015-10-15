package main

import (
	// "encoding/json"
	"fmt"
	"log"
	"net/http"

	"sync"
)

type Message struct {
	User    string
	Message string
}

var messages = make([]Message, 0)
var users = make(map[string]bool)
var lock = &sync.Mutex{}

func main() {
	http.HandleFunc("/getall", GetMessages)
	http.HandleFunc("/send", SendMessage)
	http.HandleFunc("/login", LoginUser)
	http.HandleFunc("/logout", LogoutUser)
	http.ListenAndServe(":7777", nil)
}

func LoginUser(rw http.ResponseWriter, req *http.Request) {
	lock.Lock()
	query := req.URL.Query()
	user := query.Get("user")
	log.Printf("User %s logged in", user)
	if user == "" {
		http.Error(rw, "User argument not obtained", 405)
	}
	if users[user] {
		http.Error(rw, "User already exist", 405)
	}
	users[user] = true
	rw.Write([]byte("OK"))
	lock.Unlock()
}

func LogoutUser(rw http.ResponseWriter, req *http.Request) {
	lock.Lock()
	query := req.URL.Query()
	user := query.Get("user")
	log.Printf("User %s logged out", user)
	if user == "" {
		http.Error(rw, "User argument not obtained", 405)
	}
	users[user] = false
	rw.Write([]byte("OK"))
	lock.Unlock()
}

func GetMessages(rw http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	user := query.Get("user")
	log.Printf("GetMessages from user %s", user)
	lock.Lock()
	if users[user] {
		// msg := make([]Message, len(messages))

		output := ""
		for _, item := range messages {
			output = output + fmt.Sprintf("%s : %s \n", item.User, item.Message)
		}
		// copy(msg, messages)
		// encoder := json.NewEncoder(rw)
		// encoder.Encode(msg)
		rw.Write([]byte("OK\n"))
		rw.Write([]byte(output))
	} else {
		log.Printf("%s not logged in", user)
		http.Error(rw, "Not logged in", 405)
	}
	lock.Unlock()
}

func SendMessage(rw http.ResponseWriter, req *http.Request) {
	log.Println("SendMessages")
	query := req.URL.Query()
	msg := query.Get("message")
	user := query.Get("user")

	lock.Lock()
	if users[user] {
		message := Message{
			user,
			msg,
		}
		log.Println(message)
		if len(messages) > 50 {
			messages = messages[1:]
		}
		messages = append(messages, message)
		rw.Write([]byte("OK"))
	} else {
		log.Printf("%s not logged in", user)
		http.Error(rw, "Not logged in", 405)
	}
	lock.Unlock()
}
