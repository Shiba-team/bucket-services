package middlewares

import (
	"bucket/controller"
	"bucket/model"
	"bucket/service"
	"log"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

type Options int

const (
	AUTHEN_REQUIRED Options = iota
	UN_REQUIRED
)

func Authentication(options Options) gin.HandlerFunc {
	if options == AUTHEN_REQUIRED {
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
	if options == UN_REQUIRED {
		return func(c *gin.Context) {
			token := c.Request.Header.Get("x-auth-token")
			if token == "" {
				log.Println("unAuth")
			} else {
				log.Println("Token", token)
				user, err := service.VerifyUser(token)
				if err == nil {
					log.Println("Auth")
					c.Set("user", user)
				}
			}
			c.Next()
		}
	}

	return nil
}

func FilePermission(permission model.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {

		bucketID, _ := c.Params.Get("bucketID")
		Filename, _ := c.Params.Get("filename")
		filename := path.Base(Filename)
		user := controller.GetUser(c)

		log.Print("user: ", user.Username)
		log.Print("FilePermission filename: ", filename)

		if !service.IsExistsFileInBucket(bucketID, filename) {
			c.JSON(400, gin.H{
				"error": "file not found",
			})
			c.Abort()
			return
		} else {
			file := service.GetFileInBucket(bucketID, filename)
			var bucket model.Bucket
			service.GetBucketByID(bucketID, func(err string) {}, func(b model.Bucket) {
				bucket = b
			})

			if checkBucketOwner(bucket, user.Username) {
				c.Next()
				return
			}

			if checkFileAllowAccess(file, permission) {
				log.Println("Author")
				c.Next()
				return
			}

			c.JSON(http.StatusNonAuthoritativeInfo, gin.H{
				"error": "access denied",
			})
			c.Abort()
		}
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

	log.Print("bucket status: ", bucket.Status)
	if bucket.Status == model.StatusPublic {
		return true
	}
	log.Print("check bucket status: ", "Unpublic")
	return false
}

func checkBucketOwner(bucket model.Bucket, username string) bool {
	if bucket.Owner == username {
		return true
	}

	return false
}

func checkBucketPermission(bucket model.Bucket, permission model.Permission) bool {
	if bucket.Permission == model.PermissionRead || bucket.Permission == model.PermissionReadAndWrite {
		return true
	}

	log.Print("check File status: ", "Unpublic")
	return false
}

func checkFileAllowAccess(file model.File, permission model.Permission) bool {
	if (file.Permission == permission || file.Permission == model.PermissionReadAndWrite) && file.Status == model.StatusPublic {
		return true
	}
	return false
}
