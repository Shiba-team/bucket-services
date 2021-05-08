package middlewares

import (
	"bucket/controller"
	"bucket/model"
	"bucket/service"
	"log"
	"net/http"

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

func Authorization(permission model.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {

		bucketID, _ := c.Params.Get("bucketID")
		user := controller.GetUser(c)

		log.Println(bucketID)

		if !service.GetBucketByID(bucketID, func(err string) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			c.Abort()
		}, func(bucket model.Bucket) {

			log.Printf("bucket owner: %s, username: %s", bucket.Owner, user.Username)
			if checkBucketOwner(bucket, user.Username) {
				c.Next()
				return
			}

			if checkBucketStatus(bucket) && checkBucketPermission(bucket, permission) {
				log.Println("Author")
				c.Next()
				return
			}

			c.JSON(http.StatusNonAuthoritativeInfo, gin.H{
				"error": "access denied",
			})
			c.Abort()
		}) {
			return
		}

	}
}

func checkBucketStatus(bucket model.Bucket) bool {

	if bucket.Status == model.StatusPublic {
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

func checkBucketPermission(bucket model.Bucket, permission model.Permission) bool {
	return bucket.Permission == permission || bucket.Permission == model.PermissionReadAndWrite
}
