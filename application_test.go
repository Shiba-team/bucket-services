package main

import (
	"bucket/config"
	"bucket/model"
	"bucket/service"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var testbucketid string

func TestGetAllUser(t *testing.T) {
	LoadEnv()
	config.ConnectDatabase()
	router := setupRouter()

	host := "/4"

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", host, nil)
	req.Header.Set("x-auth-token", "abc")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	type Result struct {
		Result []model.Bucket
		Count  int
	}

	var r Result
	rerr := json.Unmarshal(w.Body.Bytes(), &r)

	log.Print("length: ", len(r.Result))
	if rerr != nil {
		log.Println(rerr.Error())
		t.Error("Response type Error!")
	}

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

func TestChangePermission(t *testing.T) {
	LoadEnv()
	config.ConnectDatabase()
	router := setupRouter()

	host := os.Getenv("HOST") + "/4/" + testbucketid

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", host, nil)
	req.Header.Set("x-auth-token", "abc")
	router.ServeHTTP(w, req)

	var resultBucket model.Bucket
	rerr := json.Unmarshal(w.Body.Bytes(), &resultBucket)
	if rerr != nil {
		log.Println(rerr.Error())
		t.Error("Response type Error!")
	}
}

func TestGetBucket(t *testing.T) {
	LoadEnv()
	config.ConnectDatabase()
	router := setupRouter()

	host := os.Getenv("HOST") + "/4/" + testbucketid

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", host, nil)
	req.Header.Set("x-auth-token", "abc")
	router.ServeHTTP(w, req)

	var resultBucket model.Bucket
	rerr := json.Unmarshal(w.Body.Bytes(), &resultBucket)
	if rerr != nil {
		log.Println(rerr.Error())
		t.Error("Response type Error!")
	}

}

func TestUpdateBucket(t *testing.T) {
	LoadEnv()
	config.ConnectDatabase()
	router := setupRouter()

	host := os.Getenv("HOST") + "/4/" + testbucketid

	testBucket := model.Bucket{
		Status:     0,
		Permission: 0,
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

	scss := service.GetBucketByID(testbucketid, func(err string) {
		t.Error(err)
	}, func(bucket model.Bucket) {
		if !(bucket.Permission == 0 && bucket.Status == 0) {
			t.Error("songthing went wrong when update!")
		}
	})

	if !scss {
		t.Error("Cannot get bucket")
	}
}

func TestAddFileToBucket(t *testing.T) {
	LoadEnv()
	config.ConnectDatabase()
	router := setupRouter()

	host := os.Getenv("HOST") + "/4/" + testbucketid
	wt := httptest.NewRecorder()

	//make request
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	// Add your image file
	f, err := os.Open("./test/testfile.txt")
	if err != nil {
		return
	}

	//test
	data, _ := ioutil.ReadAll(f)

	log.Println("file size n:", len(data))

	defer f.Close()
	fw, err := w.CreateFormFile("file", "testfile.txt")
	if err != nil {
		return
	}

	if _, err = fw.Write(data); err != nil {
		return
	}
	// if _, err = io.Copy(fw, f); err != nil {
	// 	return
	// }

	// Add the other fields
	if fw, err = w.CreateFormField("filename"); err != nil {
		return
	}
	if _, err = fw.Write([]byte("testfilename.txt")); err != nil {
		return
	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	log.Println("len b: ", len(b.Bytes()))
	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", host, &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("x-auth-token", "abc")
	router.ServeHTTP(wt, req)
	assert.Equal(t, 200, wt.Code)

	log.Println("res body:", wt.Body.String())
	service.GetBucketByID(testbucketid, func(err string) {
		t.Error(err)
	}, func(bucket model.Bucket) {
		log.Print("list file: ", bucket.ListFile)
		assert.Equal(t, 1, len(bucket.ListFile))
	})

	objectId, _ := primitive.ObjectIDFromHex("60a6654c586be7718067877c")

	service.DownloadFile(objectId)

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
