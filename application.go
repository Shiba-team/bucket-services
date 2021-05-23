package main

import (
	"bucket/config"
	"bucket/controller"
	"bucket/middlewares"
	"bucket/model"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func setupRouter() *gin.Engine {

	r := gin.Default()

	client := r.Group("/4")
	{
		client.POST("", middlewares.Authentication(middlewares.AUTHEN_REQUIRED), controller.CreateBucket) // create bucket
		client.GET("", middlewares.Authentication(middlewares.AUTHEN_REQUIRED), controller.GetListBucket) // get list bucket of user
		bucket := client.Group("/:bucketID")
		{
			bucket.GET("", middlewares.Authentication(middlewares.AUTHEN_REQUIRED), middlewares.Authorization(model.PermissionRead), controller.GetBucket)                            // Get bucket infomation
			bucket.POST("", middlewares.Authentication(middlewares.AUTHEN_REQUIRED), middlewares.Authorization(model.PermissionWrite), controller.AddFileToBucket)                    // add file to bucket
			bucket.PATCH("", middlewares.Authentication(middlewares.AUTHEN_REQUIRED), middlewares.Authorization(model.PermissionWrite), controller.UpdateBucketPermissionAndStatus)   // update bucket permission and status
			bucket.DELETE("", middlewares.Authentication(middlewares.AUTHEN_REQUIRED), middlewares.Authorization(model.PermissionWrite), controller.DeleteBucket)                     // delete bucket
			bucket.GET("/1", middlewares.Authentication(middlewares.AUTHEN_REQUIRED), middlewares.Authorization(model.PermissionRead), controller.GetBucketSize)                      // get size of bucket
			bucket.GET("/2", middlewares.Authentication(middlewares.AUTHEN_REQUIRED), middlewares.Authorization(model.PermissionRead), controller.GetListFileInBucket)                // get list file of bucket
			bucket.GET("/1/*filename", middlewares.Authentication(middlewares.AUTHEN_REQUIRED), middlewares.Authorization(model.PermissionRead), controller.GetFileInfoFromBucket)    // getfile infomation
			bucket.GET("/0/*filename", middlewares.Authentication(middlewares.UN_REQUIRED), middlewares.FilePermission(model.PermissionRead), controller.GetFileFromBucket)           // download file
			bucket.PATCH("/2/*filename", middlewares.Authentication(middlewares.AUTHEN_REQUIRED), middlewares.Authorization(model.PermissionRead), controller.UpdateFilePermission)   // change file permission
			bucket.PATCH("/3/*filename", middlewares.Authentication(middlewares.AUTHEN_REQUIRED), middlewares.Authorization(model.PermissionRead), controller.UpdateFileStatus)       // change file status
			bucket.DELETE("/4/*filename", middlewares.Authentication(middlewares.AUTHEN_REQUIRED), middlewares.Authorization(model.PermissionWrite), controller.RemoveFileFromBucket) // delete file in bucket
		}

	}

	return r
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	LoadEnv()
	config.ConnectDatabase()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router := setupRouter()
	router.Run(":" + port)
}
