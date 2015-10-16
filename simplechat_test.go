package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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
		t.Errorf("Expected OK get %s", res)
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
		t.Errorf("Expected OK get %s", res)
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
		t.Errorf("Expected OK get %s", res)
	}
}
