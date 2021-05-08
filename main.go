package main

import (
	"bucket/config"
	"bucket/controller"
	"bucket/middlewares"
	"bucket/model"
	"os"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {

	r := gin.Default()

	r.Use(middlewares.Authentication())
	client := r.Group("/4")
	{
		client.POST("", controller.CreateBucket)
		client.GET("", controller.GetListBucket)
		client.POST("/:bucketID", middlewares.Authorization(model.PermissionWrite), controller.AddFileToBucket)
		client.PATCH("/:bucketID", middlewares.Authorization(model.PermissionWrite), controller.UpdateBucketPermissionAndStatus)
		client.DELETE("/:bucketID", middlewares.Authorization(model.PermissionWrite), controller.DeleteBucket)
		client.GET("/:bucketID", middlewares.Authorization(model.PermissionRead), controller.GetBucket)
		client.GET("/:bucketID/1", middlewares.Authorization(model.PermissionRead), controller.GetBucketSize)                    //get size of bucket
		client.GET("/:bucketID/2", middlewares.Authorization(model.PermissionRead), controller.GetListFileInBucket)              //get list file of bucket
		client.GET("/:bucketID/1/*filename", middlewares.Authorization(model.PermissionRead), controller.GetFileInfoFromBucket)  //getfile infomation
		client.GET("/:bucketID/0/*filename", middlewares.Authorization(model.PermissionRead), controller.GetFileFromBucket)      //download file
		client.PATCH("/:bucketID/2/*filename", middlewares.Authorization(model.PermissionRead), controller.UpdateFilePermission) //change file permission
		client.PATCH("/:bucketID/3/*filename", middlewares.Authorization(model.PermissionRead), controller.UpdateFileStatus)     //change file status

		client.DELETE("/:bucketID/*filename", middlewares.Authorization(model.PermissionWrite), controller.RemoveFileFromBucket)
	}

	return r
}

func main() {
	config.ConnectAws()
	config.ConnectDatabase()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router := setupRouter()
	router.Run(":" + port)
}
