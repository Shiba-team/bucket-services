package model

import "time"

type User struct {
	Id         string    `json:"_id"`
	Username   string    `json:"username" validate:"required"`
	Email      string    `json:"email" validate:"required"`
	Name       string    `json:"name"`
	Role       string    `json:"role" validate:"required"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	SecretKey  string    `json:"secretKey"`
}
