package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuthProvider(t *testing.T) {
	provider := NewAuthProvider()
	result := provider.isAuthenticated("user")
	if result {
		t.Errorf("Expected False got %v", result)
	}
	result = provider.login("user")
	if !result {
		t.Errorf("Expected True got %v", result)
	}
	result = provider.isAuthenticated("user")
	if !result {
		t.Errorf("Expected True got %v", result)
	}
	result = provider.logout("user")
	if !result {
		t.Errorf("Expected True got %v", result)
	}
	result = provider.isAuthenticated("user")
	if result {
		t.Errorf("Expected False got %v", result)
	}
}

func TestLoginHandler(t *testing.T) {
	requestHandler := NewRequestHandler()
	ts := httptest.NewServer(http.HandlerFunc(requestHandler.LoginUser))
	defer ts.Close()
	res, err := http.Get(ts.URL + "?user=radek")
	if err != nil {
		t.Error(err)
	}

	cont, err := ioutil.ReadAll(res.Body)
	resString := string(cont)
	defer res.Body.Close()
	if resString != "OK" {
		t.Errorf("Expected OK get %s", resString)
	}
}

func TestSendMessageHandler(t *testing.T) {
	requestHandler := NewRequestHandler()
	requestHandler.auth.login("radek")
	ts := httptest.NewServer(http.HandlerFunc(requestHandler.SendMessage))
	defer ts.Close()
	res, err := http.Get(ts.URL + "?user=radek&message=ahoj")
	if err != nil {
		t.Error(err)
	}

	cont, err := ioutil.ReadAll(res.Body)
	resString := string(cont)
	defer res.Body.Close()
	if resString != "OK" {
		t.Errorf("Expected OK get %s", resString)
	}
}

func TestGetMessagesHandler(t *testing.T) {
	requestHandler := NewRequestHandler()
	requestHandler.auth.login("radek")
	ts := httptest.NewServer(http.HandlerFunc(requestHandler.GetMessages))
	defer ts.Close()
	res, err := http.Get(ts.URL + "?user=radek")
	if err != nil {
		t.Error(err)
	}

	cont, err := ioutil.ReadAll(res.Body)
	resString := string(cont)
	defer res.Body.Close()

	if !strings.Contains(resString, "OK") {
		t.Errorf("Expected OK get %s", resString)
	}
}

func TestLogoutHandler(t *testing.T) {
	requestHandler := NewRequestHandler()
	requestHandler.auth.login("radek")
	ts := httptest.NewServer(http.HandlerFunc(requestHandler.LogoutUser))
	defer ts.Close()
	res, err := http.Get(ts.URL + "?user=radek")
	if err != nil {
		t.Error(err)
	}

	cont, err := ioutil.ReadAll(res.Body)
	resString := string(cont)
	defer res.Body.Close()

	if resString != "OK" {
		t.Errorf("Expected OK get %s", resString)
	}
}
