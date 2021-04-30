package service

import (
	"bucket/config"
	"bucket/model"
	"context"
	"time"
)

func GetFile(key string) []byte {
	// CheckFileExists(key)
	return nil
}

// func CreateFile(bucket model.Bucket, file model.File) (string, error) {

// 	if IsExistsFileName(bucket.BucketName, file.FileName) {
// 		return "", fmt.Errorf("filename %s is exists", file.FileName)
// 	}

// 	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
// 	defer cancel()

// 	_, err := config.FileCollection.InsertOne(ctx, file)
// 	if err != nil {
// 		return "", err
// 	}
// 	return file.FileName, nil
// }

func createFileInDatabase(file model.File) (string, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	_, err := config.FileCollection.InsertOne(ctx, file)
	if err != nil {
		return "", err
	}
	return file.FileName, nil
}

func IsExistsFileName(bucketname string, filename string) bool {

	lFile := GetListFileOfBucket(bucketname)

	for i := 0; i < len(lFile); i++ {
		f := lFile[i]
		if f.FileName == filename {
			return true
		}
	}

	return false
}
