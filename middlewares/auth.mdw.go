package middlewares

import (
	"bucket/controller"
	"bucket/model"
	"bucket/service"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("x-auth-token")
		if token == "" {
			log.Println("unAuth")
			c.JSON(400, gin.H{
				"error": "x-auth-token is required",
			})
			c.Abort()
			return
		}
		log.Println("Token", token)

		user, err := service.VerifyUser(token)

		if err != nil {
			log.Println("unAuth")
			c.JSON(400, gin.H{
				"error": "No Authentication",
			})
			c.Abort()
			return
		}

		log.Println("Auth")
		c.Set("user", user)
		c.Next()

	}
}

func Authorization(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {

		bucket_name, _ := c.Params.Get("bucket")
		user := controller.GetUser(c)

		bucket := service.GetBucketByName(bucket_name)

		if checkBucketOwner(bucket, user.Username) {
			c.Next()
			return
		}

		if checkBucketStatus(bucket) && checkBucketPermission(bucket, permission) {
			log.Println("Author")
			c.Next()
			return
		}

		c.JSON(403, gin.H{
			"error": "access denied",
		})
		c.Abort()
	}
}

func checkBucketStatus(bucket model.Bucket) bool {

	if bucket.Status == "public" {
		return true
	}
	return false
}

func checkBucketOwner(bucket model.Bucket, username string) bool {
	if bucket.Owner == username {
		return true
	}

	return false
}

func checkBucketPermission(bucket model.Bucket, permission string) bool {
	return strings.Contains(string(bucket.Permission), permission)
}
