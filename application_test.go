package main

import (
	"bucket/config"
	"bucket/model"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testbucketid string

func TestGetAllUser(t *testing.T) {
	LoadEnv()
	config.ConnectDatabase()
	router := setupRouter()

	host := os.Getenv("HOST") + "/4"

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", host, nil)
	req.Header.Set("x-auth-token", "abc")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	log.Println(w.Body.String())

}

func TestCreateBucket(t *testing.T) {
	LoadEnv()
	config.ConnectDatabase()
	router := setupRouter()

	host := os.Getenv("HOST") + "/4"

	testBucket := model.Bucket{
		Owner:      "testingUser",
		BucketName: "testingBucketname",
		Status:     1,
		Permission: 0,
	}

	body, err := json.Marshal(testBucket)

	if err != nil {
		t.Error("error when convert bucket object to json string")
	}

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", host, bytes.NewBuffer(body))
	req.Header.Set("x-auth-token", "abc")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var resultBucket model.Bucket

	rerr := json.Unmarshal(w.Body.Bytes(), &resultBucket)

	if rerr != nil {
		log.Println(rerr.Error())
		t.Error("Response type Error!")
	}

	if resultBucket.BucketName != testBucket.BucketName {
		t.Error("Create bucket failed!")
	} else {

		log.Println("testID: ", w.Body.String())
		testbucketid = resultBucket.ID.Hex()
		log.Print("ID: ", testbucketid)
	}
}

func TestUpdateBucket(t *testing.T) {
	LoadEnv()
	config.ConnectDatabase()
	router := setupRouter()

	host := os.Getenv("HOST") + "/4/" + testbucketid

	testBucket := model.Bucket{
		Status:     1,
		Permission: 1,
	}

	body, err := json.Marshal(testBucket)

	if err != nil {
		t.Error("error when convert bucket object to json string")
	}

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("PATCH", host, bytes.NewBuffer(body))
	req.Header.Set("x-auth-token", "abc")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestDeleteBucket(t *testing.T) {
	LoadEnv()
	config.ConnectDatabase()
	router := setupRouter()

	host := os.Getenv("HOST") + "/4/" + testbucketid

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("DELETE", host, nil)
	req.Header.Set("x-auth-token", "abc")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
