package controller

import (
	"bucket/model"
	"log"

	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) model.User {
	i, _ := c.Get("user")
	u, issuccess := i.(model.User)
	log.Println(issuccess)
	return u
}
