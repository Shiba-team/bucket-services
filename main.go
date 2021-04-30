package main

import (
	"bucket/config"
	"bucket/controller"
	"bucket/middlewares"
	"os"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {

	r := gin.Default()

	r.Use(middlewares.Authentication())
	client := r.Group("/api")
	{
		client.POST("/buckets", controller.CreateBucket)
		client.POST("/buckets/:bucket", middlewares.Authorization("write"), controller.AddFileToBucket)
		client.GET("/buckets/:bucket/*filename", middlewares.Authorization("read"), controller.GetFileFromBucket)
		client.DELETE("/buckets/:bucket/*filename", middlewares.Authorization("write"), controller.RemoveFileFromBucket)
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
