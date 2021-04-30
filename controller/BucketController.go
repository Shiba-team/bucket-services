package controller

import (
	"bucket/model"
	"bucket/service"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func CreateBucket(c *gin.Context) {
	var input model.Bucket

	user := GetUser(c)
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if service.IsExistBucketName(input.BucketName) {
		c.JSON(400, gin.H{
			"error": "bucket name is use",
		})
		return
	}

	bucket := model.NewBucket(input.BucketName, user.Username)

	resutl, err := service.CreateBucket(*bucket)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"bucket_name": resutl})
		return
	}
}

func AddFileToBucket(c *gin.Context) {
	bucket_name, isName := c.Params.Get("bucket")

	user := GetUser(c)

	if !isName {
		c.JSON(400, gin.H{
			"error": "bucketname is required",
		})
		return
	}

	if !service.IsExistBucketName(bucket_name) {
		c.JSON(400, gin.H{
			"error": "bucket name is not found",
		})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Println(err)
	}

	if service.IsExistsFileInBucket(bucket_name, header.Filename) {
		c.JSON(400, gin.H{
			"error": "filename is used",
		})
		return
	}

	//upload to S3
	path, err := service.Upload(file, header, bucket_name)

	file_name := filepath.Base(path)

	key := filepath.Join(bucket_name, file_name)

	f := model.NewFile(file_name, key, user.Username)

	service.AddFileToBucket(bucket_name, f)

	hostname := os.Getenv("HOST")

	downloadPath := hostname + "/" + "api/buckets" + "/" + bucket_name + "/" + file_name

	if err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to upload file",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"filepath": downloadPath,
	})
}

func RemoveFileFromBucket(c *gin.Context) {
	log.Print("Remove File...")
	bucket_name, _ := c.Params.Get("bucket")
	filepathname, _ := c.Params.Get("filename")

	filename := path.Base(filepathname)

	if !service.IsExistBucketName(bucket_name) {
		log.Print(" bucket not found")
		c.JSON(400, gin.H{
			"error": "bucket name is not found",
		})
		return
	}

	if !service.IsExistsFileInBucket(bucket_name, filename) {
		log.Print(" file not found")
		c.JSON(400, gin.H{
			"error": "file is not found",
		})
		return
	}

	log.Print("Removing...")
	s := service.RemoveObjectFromBucket(bucket_name, filename)
	if s {
		c.JSON(200, gin.H{
			"message": filename + " has removed",
		})
		return
	}

	c.JSON(400, gin.H{
		"error": "something when wrong!",
	})

}

func GetFileFromBucket(c *gin.Context) {
	log.Println("GetFile")

	bucket_name, _ := c.Params.Get("bucket")
	filepathname, _ := c.Params.Get("filename")

	filename := path.Base(filepathname)

	log.Print(filename)

	if !service.IsExistBucketName(bucket_name) {
		log.Print(" bucket not found")
		c.JSON(400, gin.H{
			"error": "bucket name is not found",
		})
		return
	}

	if !service.IsExistsFileInBucket(bucket_name, filename) {
		log.Print(" file not found")
		c.JSON(400, gin.H{
			"error": "file is not found",
		})
		return
	}

	key := filepath.Join(bucket_name, filepathname)

	_ = service.GetBucketByName(bucket_name)

	data := service.Download(key)

	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", http.DetectContentType(data))
	c.Header("Content-Length", fmt.Sprintf("%d", len(data)))
	c.Writer.Write(data) //the memory take up 1.2~1.7G
}

func GetListBucket(c *gin.Context) {
	user := GetUser(c)

	buckets := service.GetListBucketByUsername(user.Username)

	c.JSON(200, gin.H{
		"result": buckets,
		"count":  len(buckets),
	})
}

func GetListFileInBucket(c *gin.Context) {
	bucket_name, isName := c.Params.Get("bucket")

	if !isName {
		c.JSON(400, gin.H{
			"error": "bucketname is required",
		})
		return
	}

	if !service.IsExistBucketName(bucket_name) {
		c.JSON(400, gin.H{
			"error": "bucket name is not found",
		})
		return
	}

	files := service.GetListFileOfBucket(bucket_name)

	c.JSON(200, gin.H{
		"result": files,
		"count":  len(files),
	})

}
