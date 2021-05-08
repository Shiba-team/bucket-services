package model

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Bucket struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	Owner        string             `json:"owner" bson:"owner" validate:"required"`
	BucketName   string             `json:"bucketname" bson:"bucketname" validate:"required" binding:"required"`
	CreatedAt    time.Time          `json:"createdat" bson:"createdat" validate:"required"`
	ListFile     []File             `json:"listfile" bson:"listfile"`
	Status       Status             `json:"status" bson:"status" validate:"required"`
	Permission   Permission         `json:"permission" bson:"permission" validate:"required"`
	LastModified time.Time          `json:"lastmodified" bson:"lastmodified"  validate:"required"`
}

func NewBucket(bucket_name string, owner string) *Bucket {
	b := new(Bucket)
	b.ID = primitive.NewObjectID()
	b.Owner = owner
	b.BucketName = bucket_name
	b.CreatedAt = time.Now()
	b.Status = StatusPublic
	b.Permission = PermissionReadAndWrite
	b.LastModified = time.Now()

	return b
}

func NewFullBucket(bucket_name string, owner string, permission Permission, status Status) *Bucket {
	b := new(Bucket)
	b.ID = primitive.NewObjectID()
	b.Owner = owner
	b.BucketName = bucket_name
	b.CreatedAt = time.Now()
	b.Status = status
	b.Permission = permission
	b.LastModified = time.Now()
	return b
}

func (b Bucket) IsValidBucket(iferror func(string), ifvalid func()) bool {
	// validate := validator.New()
	// log.Println("bucket: ", b.Permission, b.Status)
	// if err := validate.Struct(b); err != nil {
	// 	log.Println("error", err)
	// 	iferror(err.Error())
	// 	return false
	// }

	if !b.Permission.IsValidPermission(func(s string) {
		iferror(s)
	}) {
		return false
	}

	if !b.Status.IsValidStatus(func(s string) {
		iferror(s)
	}) {
		return false
	}

	ifvalid()
	return true
}

func (b Bucket) GetBucketSize() int64 {
	list := b.ListFile
	var totalSize int64
	totalSize = 0
	for i := 0; i < len(list); i++ {
		totalSize += list[i].Size
	}
	return totalSize
}

func (b Bucket) GetFile(s3filename string) (File, error) {
	list := b.ListFile

	for i := 0; i < len(list); i++ {
		if list[i].S3Name == s3filename {
			return list[i], nil
		}
	}

	return File{}, errors.New("no file in bucket")
}

// func GetBucketCollection() *mongo.Collection {
// 	return config.Client.Database("storage").Collection("Bucket")
// }
