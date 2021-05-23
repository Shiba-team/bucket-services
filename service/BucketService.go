package service

import (
	"bucket/config"
	"bucket/model"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//return iferror
func CreateBucket(bucket model.Bucket, iferror func(err string), ifsuccess func(id string, name string)) bool {
	if !IsExistBucketName(bucket.BucketName) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		log.Println("bucket info: ", bucket)

		result, err := config.BucketCollection.InsertOne(ctx, bucket)
		if err != nil {
			iferror(err.Error())
			return false
		}
		id := result.InsertedID.(primitive.ObjectID)
		ifsuccess(id.Hex(), bucket.BucketName)
		return true
	} else {
		iferror("exists bucket name")
		return false
	}
}

func UpdateBucketPermission(bucketId string, permission model.Permission) error {
	var err error
	if !GetBucketByID(bucketId, func(errid string) {
		err = errors.New(errid)
	}, func(bucket model.Bucket) {
		newBucket := bucket
		newBucket.Permission = permission
		newBucket.LastModified = time.Now()

		log.Println("permission:", permission)
		log.Println(newBucket)
		UpdateBucket(newBucket, func(errid string) {
			err = errors.New(errid)
		}, func(newbucket model.Bucket) {
			log.Println(newbucket)
		})
	}) {
		return err
	}
	return nil
}

func UpdateBucketStatus(bucketId string, status model.Status) error {
	var err error
	if !GetBucketByID(bucketId, func(errid string) {
		err = errors.New(errid)
	}, func(bucket model.Bucket) {
		newBucket := bucket
		newBucket.Status = status
		newBucket.LastModified = time.Now()

		log.Println("status:", status)
		log.Println(newBucket)
		UpdateBucket(newBucket, func(errid string) {
			err = errors.New(errid)
		}, func(newbucket model.Bucket) {
			log.Println(newbucket)
		})
	}) {
		return err
	}
	return nil
}

// func UpdateBucket(bucket model.Bucket, iferror func(err string), ifsuccess func(newbucket model.Bucket)) bool {

// }

func DeleteBucket(ID string, iferror func(err string), ifsuccess func(name string)) {
	GetBucketByID(ID, func(err string) {
		iferror(err)
	}, func(bucket model.Bucket) {
		err := deleteFileInBucket(bucket.ListFile)
		if err != nil {
			log.Print("DeleteBucket: error when delete file in bucket: ", err)
		} else {
			deleteBucket(bucket.ID, func(err string) {
				iferror(err)
			}, func(name string) {
				ifsuccess(name)
			})
		}

	})
}

func deleteFileInBucket(files []model.File) error {
	for i := 0; i < len(files); i++ {
		FileID, err := primitive.ObjectIDFromHex(files[i].FileID)
		if err != nil {
			log.Print("deleteFileBucket: ", err)
			return err
		}
		DeleteFile(FileID)
	}
	return nil
}

func deleteBucket(ID primitive.ObjectID, iferror func(err string), ifsuccess func(name string)) bool {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var bucket model.Bucket

	result := config.BucketCollection.FindOneAndDelete(ctx, bson.M{"_id": ID})

	err := result.Decode(&bucket)

	if err != nil {
		iferror(err.Error())
		return false
	}

	ifsuccess(bucket.BucketName)
	return true
}

func IsExistBucketName(bucketname string) bool {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var bucket model.Bucket

	result := config.BucketCollection.FindOne(ctx, bson.M{"bucketname": bucketname})

	err := result.Decode(&bucket)

	if err != nil {
		log.Println("IsExistsBucketName: ", err)
	}

	return bucket.BucketName == bucketname
}

func GetBucketByName(bucketname string) model.Bucket {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var bucket model.Bucket

	result := config.BucketCollection.FindOne(ctx, bson.M{"bucketname": bucketname})

	err := result.Decode(&bucket)

	if err != nil {
		log.Println("GetBucketByName: ", err)
	}

	return bucket
}

func GetBucketByID(id string, iferror func(err string), ifsuccess func(bucket model.Bucket)) bool {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var bucket model.Bucket

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		iferror(err.Error())
		return false
	}

	result := config.BucketCollection.FindOne(ctx, bson.M{"_id": objectId})

	err = result.Decode(&bucket)

	log.Println("err", err)
	// log.Println(bucket)
	if err != nil {
		iferror(err.Error())
		return false
	}

	ifsuccess(bucket)
	return true
}

//Get list of bucket via bucket name
func GetListFileOfBucket(id string) []model.File {
	var list []model.File

	GetBucketByID(id, func(err string) {

	}, func(bucket model.Bucket) {
		list = bucket.ListFile
	})

	return list
}

// add file to bucket
func AddFileToBucket(bucketname string, file model.File) (string, error) {

	bucket, _, _ := addFileToDatabase(bucketname, file)

	var err error
	UpdateBucket(*bucket, func(errs string) { err = errors.New(errs) }, func(newbucket model.Bucket) {})
	if err != nil {
		return "", err
	}
	return bucket.BucketName + "/" + file.FileName, nil
}

//Add file to database in cloude
func addFileToDatabase(bucketID string, file model.File) (*model.Bucket, string, error) {
	var bucketf model.Bucket
	var erro error
	if GetBucketByID(bucketID, func(err string) {
		erro = errors.New(err)
	}, func(bucket model.Bucket) {
		bucket.ListFile = append(bucket.ListFile, file)
		bucketf = bucket
	}) {
		return &bucketf, file.FileName, nil
	} else {
		return nil, "", erro
	}

	//filename, _ := createFileInDatabase(file)

	// if filename == "" {
	// 	return nil, "", fmt.Errorf("can't create file")
	// }

	// file.FileName = filename

}

func UpdateBucket(bucket model.Bucket, iferror func(err string), ifsuccess func(newbucket model.Bucket)) bool {

	return bucket.IsValidBucket(func(s string) {
		iferror(s)
	}, func() {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		_, err := config.BucketCollection.ReplaceOne(ctx,
			bson.M{"_id": bucket.ID},
			bucket,
		)

		if err != nil {
			iferror(err.Error())
		}

		ifsuccess(bucket)
	})

}

func RemoveObjectFromBucket(bucketID string, filename string) bool {
	log.Print("Removing in database")
	ok, id, _ := deleteFileInDatabase(bucketID, filename)
	if ok {
		fileID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return false
		}
		ferr := DeleteFile(fileID)
		if ferr != nil {
			log.Print("Remove Object ", ferr)
			return false
		}
	}
	return ok
}

func deleteFileInDatabase(bucketID string, filename string) (bool, string, error) {

	log.Print("Removing in database...")
	var bucket model.Bucket
	var er error
	GetBucketByID(bucketID, func(err string) {
		er = errors.New(err)
	}, func(rbucket model.Bucket) {
		bucket = rbucket
	})

	if er != nil {
		return false, "", er
	}

	pos := findFilePosition(bucket.ListFile, filename)
	if pos == -1 {
		log.Print("file not found")
		return false, "", fmt.Errorf("unable to find filename")
	}

	deletedFileID := bucket.ListFile[pos].FileID
	bucket.ListFile = append(bucket.ListFile[:pos], bucket.ListFile[pos+1:]...)

	UpdateBucket(bucket, func(err string) {}, func(newbucket model.Bucket) {})
	log.Print("Deleted in database")
	return true, deletedFileID, nil
}

func findFilePosition(list []model.File, s3filename string) int {
	for i := 0; i < len(list); i++ {
		if list[i].S3Name == s3filename {
			return i
		}
	}
	return -1
}

func IsExistsFileInBucket(id string, s3filename string) bool {
	list := GetListFileOfBucket(id)

	pos := findFilePosition(list, s3filename)
	return pos != -1
}

func GetListBucketByUsername(username string) []model.Bucket {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var buckets []model.Bucket

	result, error := config.BucketCollection.Find(ctx, bson.M{"owner": username})

	if error != nil {
		return nil
	}

	err := result.All(ctx, &buckets)

	if err != nil {
		log.Println("GetBucketByName: ", err)
		return nil
	}

	return buckets
}

func IsExistsBucketID(id string) bool {
	return GetBucketByID(id, func(err string) {
		log.Println(err)
	}, func(bucket model.Bucket) {})
}
