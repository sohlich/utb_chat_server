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

type Authenticator interface {
	isAuthenticated(user string) bool
	login(user string) bool
	logout(user string) bool
}

type AuthProvider struct {
	users    map[string]bool
	userLock *sync.Mutex
}

func NewAuthProvider() *AuthProvider {
	aProv := &AuthProvider{
		make(map[string]bool),
		&sync.Mutex{},
	}
	return aProv
}

func (a *AuthProvider) isAuthenticated(user string) bool {
	isAuth := false
	a.userLock.Lock()
	isAuth = a.users[user]
	a.userLock.Unlock()
	return isAuth
}

func (a *AuthProvider) login(user string) bool {
	isLog := false
	a.userLock.Lock()
	exist := a.users[user]
	if !exist {
		a.users[user] = true
		isLog = true
	}
	a.userLock.Unlock()
	return isLog
}

func (a *AuthProvider) logout(user string) bool {
	a.userLock.Lock()
	a.users[user] = false
	a.userLock.Unlock()
	return true
}

type RequestHandler struct {
	auth        Authenticator
	messages    []Message
	messageLock *sync.Mutex
}

func NewRequestHandler() *RequestHandler {
	authHandlr := &RequestHandler{
		NewAuthProvider(),
		make([]Message, 0),
		&sync.Mutex{},
	}
	return authHandlr
}

func (a *RequestHandler) LoginUser(rw http.ResponseWriter, req *http.Request) {

	query := req.URL.Query()
	user := query.Get("user")
	log.Printf("User %s logged in", user)
	if user == "" {
		http.Error(rw, "User argument not obtained", 405)
		return
	}
	if !a.auth.login(user) {
		http.Error(rw, "User already exist", 405)
		return
	}
	rw.Write([]byte("OK"))
}

func (a *RequestHandler) LogoutUser(rw http.ResponseWriter, req *http.Request) {

	query := req.URL.Query()
	user := query.Get("user")
	log.Printf("User %s logged out", user)
	if user == "" {
		http.Error(rw, "User argument not obtained", 405)
		return
	}
	a.auth.logout(user)
	rw.Write([]byte("OK"))
}

func (a *RequestHandler) GetMessages(rw http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	user := query.Get("user")
	log.Printf("GetMessages from user %s", user)

	if a.auth.isAuthenticated(user) {

		a.messageLock.Lock()
		output := ""
		for _, item := range a.messages {
			output = output + fmt.Sprintf("%s : %s \n", item.User, item.Message)
		}
		a.messageLock.Unlock()
		rw.Write([]byte("OK\n"))
		rw.Write([]byte(output))
	} else {
		log.Printf("%s not logged in", user)
		http.Error(rw, "Not logged in", 405)
	}
}

func (a *RequestHandler) SendMessage(rw http.ResponseWriter, req *http.Request) {
	log.Println("SendMessages")
	query := req.URL.Query()
	msg := query.Get("message")
	user := query.Get("user")

	if a.auth.isAuthenticated(user) {
		message := Message{
			user,
			msg,
		}
		log.Println(message)
		a.messageLock.Lock()
		if len(a.messages) > 50 {
			a.messages = a.messages[1:]
		}
		a.messages = append(a.messages, message)
		a.messageLock.Unlock()
		rw.Write([]byte("OK"))
	} else {
		log.Printf("%s not logged in", user)
		http.Error(rw, "Not logged in", 405)
	}
}

func main() {
	authHndlr := NewRequestHandler()
	http.HandleFunc("/getall", authHndlr.GetMessages)
	http.HandleFunc("/send", authHndlr.SendMessage)
	http.HandleFunc("/login", authHndlr.LoginUser)
	http.HandleFunc("/logout", authHndlr.LogoutUser)
	http.ListenAndServe(":7777", nil)
}
