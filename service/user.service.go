package service

import "bucket/model"

func VerifyUser(token string) (model.User, error) {
	var user model.User

	user.Username = "testuser"
	return user, nil
}
