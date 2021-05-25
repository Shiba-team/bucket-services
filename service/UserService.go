package service

import (
	"bucket/model"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	// "os"
)

func sendReq(header http.Header) *http.Response {
	userhost := os.Getenv("USER_HOST")
	url := userhost + "/api/auth/verify-token"
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	req.Header = header

	res, err := client.Do(req)
	if err != nil {
		log.Print("sendReq: error when send request: ", err)
		return nil
	}
	return res
}

func VerifyUser(header http.Header) (model.User, error) {
	res := sendReq(header)
	if res == nil {
		log.Print("VerifyUser: Can't send requrest!")
		return model.User{}, errors.New("can not send request")
	}
	defer res.Body.Close()

	var user model.User
	//check body
	log.Print("VerifyUser: responsebody: ", res.Body)

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Print("VerifyUser: Error when Read response body: ", err)
		return model.User{}, err
	}
	bodyString := string(bodyBytes)
	log.Print("VerifyUser: bodystring", bodyString)
	rerr := json.Unmarshal(bodyBytes, &user)
	if rerr != nil {
		log.Println("VerifyUser: Error when convert body to User object: ", rerr.Error())
		return model.User{}, rerr
	}

	log.Print("VerifyUser: username: ", user.Username)

	return user, nil
}
