package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
	"user-management-api/model"
)

func TestCreateUserController(t *testing.T) {
	userInfo := model.User{
		Email:     "test@gmail.com",
		Fullname:  "Mr. Test",
		Telephone: "15200",
		Password:  "123456",
	}
	userRgFromJSON, _ := json.Marshal(userInfo)
	req, _ := http.NewRequest("POST", "http://localhost:8080/api/registration", bytes.NewBuffer(userRgFromJSON))
	var responseObject model.UserObject
	client := &http.Client{Timeout: time.Second * 60}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &responseObject)
	want := "test@gmail.com"
	if responseObject.Email != want {
		t.Errorf("got %s want %s", responseObject.Email, want)
	}
}
