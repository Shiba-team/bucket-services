package service

import (
	"bucket/config"
	"bucket/model"
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
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

// func createFileInDatabase(file model.File) (string, error) {
// 	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
// 	defer cancel()

// 	_, err := config.FileCollection.InsertOne(ctx, file)
// 	if err != nil {
// 		return "", err
// 	}
// 	return file.FileName, nil
// }

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

func UploadFile(data []byte, filename string) string {

	// data, err := ioutil.ReadFile(file)
	// if err != nil {
	//     log.Fatal(err)
	// }
	conn := config.InitiateMongoClient()
	bucket, err := gridfs.NewBucket(
		conn.Database("myfiles"),
	)
	if err != nil {
		log.Print("Upload File: ", err)
		return ""
	}

	FileID := primitive.NewObjectID()

	uploadStream, err := bucket.OpenUploadStreamWithID(FileID, filename)
	// uploadStream, err := bucket.OpenUploadStream(
	// 	filename,
	// )
	if err != nil {
		log.Print("error when upload: ", err)
	}
	defer uploadStream.Close()

	fileSize, err := uploadStream.Write(data)
	if err != nil {
		log.Print("Upload file: ", "Write data: ", err)
	}
	log.Printf("Write file to DB was successful. File %s size: %d M\n", FileID.Hex(), fileSize)
	return FileID.Hex()
}
func DownloadFile(ID primitive.ObjectID) []byte {
	conn := config.InitiateMongoClient()

	// For CRUD operations, here is an example
	db := conn.Database("myfiles")
	fsFiles := db.Collection("fs.files")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	var results bson.M
	err := fsFiles.FindOne(ctx, bson.M{}).Decode(&results)
	if err != nil {
		log.Fatal(err)
	}
	// you can print out the results
	fmt.Println("result search: ", results)

	bucket, _ := gridfs.NewBucket(
		db,
	)
	var buf bytes.Buffer
	dStream, err := bucket.DownloadToStream(ID, &buf)
	// dStream, err := bucket.DownloadToStreamByName(fileName, &buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("File size to download: %v\n", dStream)
	// ioutil.WriteFile("name.txt", buf.Bytes(), 0600)
	return buf.Bytes()

}

func DeleteFile(ID primitive.ObjectID) error {

	conn := config.InitiateMongoClient()

	// check file
	db := conn.Database("myfiles")
	fsFiles := db.Collection("fs.files")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var results bson.M
	err := fsFiles.FindOne(ctx, bson.M{}).Decode(&results)
	if err != nil {
		log.Print("error when search file: ", err)
	}
	// you can print out the results
	fmt.Println(results)

	bucket, err := gridfs.NewBucket(
		db,
	)
	if err != nil {
		log.Print("error when get bucket: ", err)
		os.Exit(1)
	}

	log.Println("File ID: ", ID)
	return bucket.Delete(ID)
}

func GetFileInBucket(bucketID, filename string) model.File {
	files := GetListFileOfBucket(bucketID)
	for i := 0; i < len(files); i++ {
		if files[i].S3Name == filename {
			return files[i]
		}
	}
	return model.File{}
}
