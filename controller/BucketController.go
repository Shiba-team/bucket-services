package controller

import (
	"bucket/model"
	"bucket/service"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateBucket(c *gin.Context) {

	ifError := func(err string) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
	ifGetBucketInputSuccess := func(b model.Bucket) {
		b.IsValidBucket(ifError, func() {
			if service.IsExistBucketName(b.BucketName) {
				c.JSON(400, gin.H{
					"error": "bucket name is used",
				})
				return
			}
			service.CreateBucket(b, ifError, func(id, name string) {
				c.JSON(http.StatusOK, gin.H{
					"bucketname": name,
					"_id":        id,
				})
			})
		})
	}

	getBucketInput(c, ifError, ifGetBucketInputSuccess)

}

func UpdateBucket(c *gin.Context) {

	getBucketInput(c, func(err string) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}, func(b model.Bucket) {

	})

}

func UpdateBucketPermissionAndStatus(c *gin.Context) {

	bucketID, ok := c.Params.Get("bucketID")

	if !ok {
		return
	}

	var f interface{}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to read request body",
		})
	}
	input, ok := f.(map[string]interface{})

	if ok {
		log.Println(input)

		p, exists := input["permission"]

		if exists {
			p64 := p.(float64)
			permission := model.Permission(p64)
			log.Println("per col:", permission)
			service.UpdateBucketPermission(bucketID, permission)
		}

		s, exists := input["status"]

		if exists {
			s64 := s.(float64)
			status := model.Status(s64)
			service.UpdateBucketStatus(bucketID, status)
		} else {
			log.Println("do not have status")
		}

		c.JSON(200, gin.H{
			"message": "Updated",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad body",
		})
	}

}

func getBucketInput(c *gin.Context, iferror func(err string), ifsuccess func(model.Bucket)) bool {
	var input model.Bucket

	if err := c.ShouldBindJSON(&input); err != nil {
		iferror(err.Error())
		return false
	}

	user := GetUser(c)

	if input.BucketName != "" {
		if !input.Permission.IsValidPermission(func(s string) {
			iferror(s)
		}) {
			return false
		}

		if !input.Status.IsValidStatus(func(s string) {
			iferror(s)
		}) {
			return false
		}

		bucket := model.NewFullBucket(input.BucketName, user.Username, input.Permission, input.Status)

		return bucket.IsValidBucket(func(s string) {
			log.Println("error invalid bucket")
			iferror(s)
		}, func() {
			ifsuccess(*bucket)
		})
	}

	iferror("BucketName is required")
	return false
}

func GetBucketSize(c *gin.Context) {

	ID, _ := c.Params.Get("bucketID")

	service.GetBucketByID(ID, func(err string) {
		c.JSON(400, gin.H{
			"error": err,
		})
	}, func(bucket model.Bucket) {
		size := bucket.GetBucketSize()
		c.JSON(200, gin.H{
			"bucketsize": size,
		})
	})
}

func AddFileToBucket(c *gin.Context) {
	ID, _ := c.Params.Get("bucketID")
	user := GetUser(c)

	// if !isName {
	// 	c.JSON(400, gin.H{
	// 		"error": "bucketname is required",
	// 	})
	// 	return
	// }

	// if !service.IsExistsBucketID(ID) {
	// 	c.JSON(400, gin.H{
	// 		"error": "bucketID is not found",
	// 	})
	// 	return
	// }

	file, header, err := c.Request.FormFile("file")
	filename := c.Request.FormValue("filename")
	if filename == "" {
		filename = header.Filename
	}

	if err != nil {
		log.Println(err)
	}

	// if service.IsExistsFileInBucket(ID, filename) {
	// 	c.JSON(400, gin.H{
	// 		"error": "filename is used",
	// 	})
	// 	return
	// }
	s3filename, err := model.RandomHex(64)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err,
		})
	}
	s3filename += path.Ext(filename)

	// upload to S3
	path, err := service.Upload(file, header, ID, s3filename)
	//---------------------------------
	if err != nil {
		log.Println("error", err)
		c.JSON(500, gin.H{
			"error": "Failed to upload file",
		})
		return
	}
	//---------------------------------
	log.Println("path: ", path)
	// key := filepath.Join(ID, file_name)
	hostname := os.Getenv("HOST")
	downloadPath := hostname + "/4/" + ID + "/0/" + s3filename
	f := model.File{
		FileName:     filename,
		FileLink:     path,
		S3Name:       s3filename,
		DownloadLink: downloadPath,
		CreatedAt:    time.Now(),
		CreatedBy:    user.Username,
		LastModify:   time.Now(),
		Status:       model.StatusPrivate,
		Permission:   model.PermissionReadAndWrite,
		Size:         header.Size,
	}
	service.AddFileToBucket(ID, f)

	log.Println(f)

	c.JSON(http.StatusOK, gin.H{
		"filepath": downloadPath,
	})
}

func RemoveFileFromBucket(c *gin.Context) {
	log.Print("Remove File...")
	bucket_name, _ := c.Params.Get("bucketID")
	filepathname, _ := c.Params.Get("filename")

	s3filename := path.Base(filepathname)

	if !service.IsExistBucketName(bucket_name) {
		log.Print(" bucket not found")
		c.JSON(400, gin.H{
			"error": "bucket name is not found",
		})
		return
	}

	if !service.IsExistsFileInBucket(bucket_name, s3filename) {
		log.Print(" file not found")
		c.JSON(400, gin.H{
			"error": "file is not found",
		})
		return
	}

	log.Print("Removing...")
	s := service.RemoveObjectFromBucket(bucket_name, s3filename)
	if s {
		c.JSON(200, gin.H{
			"message": s3filename + " has removed",
		})
		return
	}

	c.JSON(400, gin.H{
		"error": "something when wrong!",
	})

}

func GetFileFromBucket(c *gin.Context) {
	log.Println("GetFile")

	ID, _ := c.Params.Get("bucketID")
	filepathname, _ := c.Params.Get("filename")

	s3filename := path.Base(filepathname)

	log.Print(s3filename)

	// if !service.IsExistsBucketID(ID) {
	// 	log.Print(" bucket not found")
	// 	c.JSON(400, gin.H{
	// 		"error": "bucket name is not found",
	// 	})
	// 	return
	// }

	if !service.IsExistsFileInBucket(ID, s3filename) {
		log.Print(" file not found")
		c.JSON(400, gin.H{
			"error": "file is not found",
		})
		return
	}

	service.GetBucketByID(ID, func(err string) {
		c.JSON(400, gin.H{
			"error": err,
		})
	}, func(bucket model.Bucket) {
		file, err := bucket.GetFile(s3filename)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err,
			})
		}

		// _ = service.GetBucketByName(ID)

		data := service.Download(file.FileLink)

		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.FileName))
		c.Header("Content-Type", http.DetectContentType(data))
		c.Header("Content-Length", fmt.Sprintf("%d", len(data)))
		c.Writer.Write(data) //the memory take up 1.2~1.7G
	})

}

func GetListBucket(c *gin.Context) {
	user := GetUser(c)

	buckets := service.GetListBucketByUsername(user.Username)

	c.JSON(200, gin.H{
		"result": buckets,
		"count":  len(buckets),
	})
}

func GetFileInfoFromBucket(c *gin.Context) {
	ID, _ := c.Params.Get("bucketID")
	filepathname, _ := c.Params.Get("filename")

	s3filename := path.Base(filepathname)

	var file model.File
	if service.GetBucketByID(ID, func(err string) {
		c.JSON(400, gin.H{
			"error": err,
		})
	}, func(bucket model.Bucket) {
		f, err := bucket.GetFile(s3filename)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err,
			})
		}

		file = f
	}) {
		c.JSON(200, file)
	}
}

func GetBucket(c *gin.Context) {
	bucketID, _ := c.Params.Get("bucketID")

	service.GetBucketByID(bucketID, func(err string) {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
	}, func(bucket model.Bucket) {
		c.JSON(200, bucket)
	})
}

func DeleteBucket(c *gin.Context) {
	bucketID, _ := c.Params.Get("bucketID")

	service.DeleteBucket(bucketID, func(err string) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}, func(name string) {
		c.JSON(http.StatusOK, gin.H{
			"message": name + " deleted successfully",
		})
	})
}

func GetListFileInBucket(c *gin.Context) {
	bucketID, isName := c.Params.Get("bucketID")

	if !isName {
		c.JSON(400, gin.H{
			"error": "bucketname is required",
		})
		return
	}

	if !service.IsExistsBucketID(bucketID) {
		c.JSON(400, gin.H{
			"error": "bucket name is not found",
		})
		return
	}

	files := service.GetListFileOfBucket(bucketID)

	c.JSON(200, gin.H{
		"result": files,
		"count":  len(files),
	})

}

func UpdateFilePermission(c *gin.Context) {
	ID, _ := c.Params.Get("bucketID")
	filepathname, _ := c.Params.Get("filename")

	s3filename := path.Base(filepathname)

	var f interface{}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to read request body",
		})
	}
	input, ok := f.(map[string]interface{})

	if ok {
		log.Println(input)

		p, exists := input["permission"]

		if exists {
			p64 := p.(float64)
			permission := model.Permission(p64)
			log.Println("per col:", permission)

			if !permission.IsValidPermission(func(s string) {
				c.JSON(400, gin.H{
					"error": s,
				})
			}) {
				return
			}

			var bucket model.Bucket

			if service.GetBucketByID(ID, func(err string) {
				c.JSON(400, gin.H{
					"error": err,
				})
			}, func(bucketf model.Bucket) {

				bucket = bucketf

				for i := 0; i < len(bucket.ListFile); i++ {
					if bucket.ListFile[i].S3Name == s3filename {
						bucket.ListFile[i].Permission = permission
					}
				}
			}) {
				service.UpdateBucket(bucket, func(err string) {}, func(newbucket model.Bucket) {
					c.JSON(200, gin.H{
						"message": "permission updated successfully",
					})
				})

			}
		} else {
			c.JSON(400, gin.H{
				"error": "permission is required",
			})
		}
	}

}

func UpdateFileStatus(c *gin.Context) {
	ID, _ := c.Params.Get("bucketID")
	filepathname, _ := c.Params.Get("filename")

	s3filename := path.Base(filepathname)

	var f interface{}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to read request body",
		})
	}
	input, ok := f.(map[string]interface{})

	if ok {
		log.Println(input)

		s, exists := input["status"]

		if exists {
			s64 := s.(float64)
			status := model.Status(s64)
			log.Println("per col:", status)

			if !status.IsValidStatus(func(s string) {
				c.JSON(400, gin.H{
					"error": s,
				})
			}) {
				return
			}

			var bucket model.Bucket

			if service.GetBucketByID(ID, func(err string) {
				c.JSON(400, gin.H{
					"error": err,
				})
			}, func(bucketf model.Bucket) {

				bucket = bucketf

				for i := 0; i < len(bucket.ListFile); i++ {
					if bucket.ListFile[i].S3Name == s3filename {
						bucket.ListFile[i].Status = status
					}
				}
			}) {
				service.UpdateBucket(bucket, func(err string) {}, func(newbucket model.Bucket) {
					c.JSON(200, gin.H{
						"message": "status updated successfully",
					})
				})

			}
		} else {
			c.JSON(400, gin.H{
				"error": "status is required",
			})
		}
	}
}
