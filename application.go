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

	r.Use(middlewares.Authentication())
	client := r.Group("/4")
	{
		client.POST("", controller.CreateBucket)                                                                                 // create bucket
		client.GET("", controller.GetListBucket)                                                                                 // get list bucket of user
		client.POST("/:bucketID", middlewares.Authorization(model.PermissionWrite), controller.AddFileToBucket)                  // add file to bucket
		client.PATCH("/:bucketID", middlewares.Authorization(model.PermissionWrite), controller.UpdateBucketPermissionAndStatus) // update bucket permission and status
		client.DELETE("/:bucketID", middlewares.Authorization(model.PermissionWrite), controller.DeleteBucket)                   // delete bucket
		client.GET("/:bucketID", middlewares.Authorization(model.PermissionRead), controller.GetBucket)                          // Get bucket infomation
		client.GET("/:bucketID/1", middlewares.Authorization(model.PermissionRead), controller.GetBucketSize)                    // get size of bucket
		client.GET("/:bucketID/2", middlewares.Authorization(model.PermissionRead), controller.GetListFileInBucket)              // get list file of bucket
		client.GET("/:bucketID/1/*filename", middlewares.Authorization(model.PermissionRead), controller.GetFileInfoFromBucket)  // getfile infomation
		client.GET("/:bucketID/0/*filename", middlewares.Authorization(model.PermissionRead), controller.GetFileFromBucket)      // download file
		client.PATCH("/:bucketID/2/*filename", middlewares.Authorization(model.PermissionRead), controller.UpdateFilePermission) // change file permission
		client.PATCH("/:bucketID/3/*filename", middlewares.Authorization(model.PermissionRead), controller.UpdateFileStatus)     // change file status
		client.DELETE("/:bucketID/*filename", middlewares.Authorization(model.PermissionWrite), controller.RemoveFileFromBucket) // delete file in bucket
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
	config.ConnectAws()
	config.ConnectDatabase()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router := setupRouter()
	router.Run(":" + port)
}
