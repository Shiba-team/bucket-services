package service

import (
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func Upload(file multipart.File, header *multipart.FileHeader, bucket_name string) (string, error) {
	sess, _ := session.NewSession()
	uploader := s3manager.NewUploader(sess)

	filename := header.Filename

	size := header.Size
	buffer := make([]byte, size)
	file.Read(buffer)
	//upload to the s3 bucket
	filepath, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:               aws.String(os.Getenv("AWS_BUCKET")),
		Key:                  aws.String(bucket_name + "/" + filename),
		ACL:                  aws.String("public-read"), // could be private if you want it to be access by only authorized users
		Body:                 bytes.NewReader(buffer),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
		StorageClass:         aws.String("INTELLIGENT_TIERING"),
	})

	if err != nil {
		return "", err
	}

	return filepath.Location, nil
}

func Download(key string) []byte {
	// The session the S3 Downloader will use
	sess := session.Must(session.NewSession())

	// Create a downloader with the session and default options
	downloader := s3manager.NewDownloader(sess)

	// Create a file to write the S3 Object contents to.

	f, err := os.CreateTemp("", "temp")

	if err != nil {
		log.Fatal(err)
	}

	// Write the contents of S3 Object to the file
	n, err := downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(key),
	})
	if err != nil {
		fmt.Printf("failed to download file, %v", err)
	}
	fmt.Printf("file downloaded, %d bytes\n", n)

	data := make([]byte, n)
	f.Read(data)

	return data
}

func deleteObject(bucket string, key string) (bool, error) {
	log.Print("Deleting in S3...")
	sess, _ := session.NewSession()
	svc := s3.New(sess)
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(bucket + "/" + key),
	}

	result, err := svc.DeleteObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return false, err
	}

	log.Printf("Object has been deleted")
	log.Println(result)
	return true, nil
}
