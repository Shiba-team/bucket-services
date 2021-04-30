package service

import (
	"bucket/config"
	"bucket/model"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func CreateBucket(bucket model.Bucket) (string, error) {
	if !IsExistBucketName(bucket.BucketName) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		_, err := config.BucketCollection.InsertOne(ctx, bucket)
		if err != nil {
			return "", err
		}
		return bucket.BucketName, nil
	} else {
		return "", fmt.Errorf("exists bucket name")
	}
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

//Get list of bucket via bucket name
func GetListFileOfBucket(bucketname string) []model.File {
	bucket := GetBucketByName(bucketname)

	return bucket.ListFile
}

// add file to bucket
func AddFileToBucket(bucketname string, file model.File) (string, error) {

	bucket, _, _ := addFileToDatabase(bucketname, file)

	UpdateBucket(*bucket)
	return bucket.BucketName + "/" + file.FileName, nil
}

//Add file to database in cloude
func addFileToDatabase(bucketname string, file model.File) (*model.Bucket, string, error) {
	bucket := GetBucketByName(bucketname)

	filename, _ := createFileInDatabase(file)

	if filename == "" {
		return nil, "", fmt.Errorf("can't create file")
	}

	file.FileName = filename
	bucket.ListFile = append(bucket.ListFile, file)

	return &bucket, filename, nil
}

func UpdateBucket(bucket model.Bucket) model.Bucket {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	_, err := config.BucketCollection.ReplaceOne(ctx,
		bson.M{"bucketname": bucket.BucketName},
		bucket,
	)

	if err != nil {
		log.Println("GetBucketByName: ", err)
	}

	return bucket
}

func RemoveObjectFromBucket(bucketname string, filename string) bool {
	log.Print("Removing in database")
	ok, _ := deleteFileInDatabase(bucketname, filename)
	if ok {
		log.Print(" Removing in s3")
		deleteObject(bucketname, filename)
	}
	return ok
}

func deleteFileInDatabase(bucketname string, filename string) (bool, error) {

	log.Print("Removing in database...")
	bucket := GetBucketByName(bucketname)

	pos := findFilePosition(bucket.ListFile, filename)
	if pos == -1 {
		log.Print("file not found")
		return false, fmt.Errorf("unable to find filename")
	}

	bucket.ListFile = append(bucket.ListFile[:pos], bucket.ListFile[pos+1:]...)

	UpdateBucket(bucket)
	log.Print("Deleted in database")
	return true, nil
}

func findFilePosition(list []model.File, filename string) int {
	for i := 0; i < len(list); i++ {
		if list[i].FileName == filename {
			return i
		}
	}
	return -1
}

func IsExistsFileInBucket(bucketname string, filename string) bool {
	list := GetListFileOfBucket(bucketname)

	pos := findFilePosition(list, filename)
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
