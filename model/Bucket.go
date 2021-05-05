package model

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
)

type Bucket struct {
	Id           string     `json:"_id"`
	Owner        string     `json:"owner" validate:"required"`
	BucketName   string     `json:"bucketname" validate:"required" binding:"required"`
	CreatedAt    time.Time  `json:"createdat" validate:"required"`
	ListFile     []File     `json:"listfile"`
	Status       Status     `json:"status"  validate:"required"`
	Permission   Permission `json:"permission" validate:"required"`
	LastModified time.Time  `json:"lastmodified"  validate:"required"`
}

func NewBucket(bucket_name string, owner string) *Bucket {
	b := new(Bucket)
	b.Owner = owner
	b.BucketName = bucket_name
	b.CreatedAt = time.Now()
	b.Status = "private"
	b.Permission = "read:write"
	b.LastModified = time.Now()

	return b
}

func NewFullBucket(bucket_name string, owner string, permission string, status string) *Bucket {
	b := new(Bucket)
	b.Owner = owner
	b.BucketName = bucket_name
	b.CreatedAt = time.Now()
	b.Status = Status(status)
	b.Permission = Permission(permission)
	b.LastModified = time.Now()
	return b
}

func (b Bucket) IsValidBucket() error {

	validate := validator.New()
	err := validate.Struct(b)

	if err := b.Permission.IsValidPermission(); err != nil {
		return errors.New("invalid permission")
	}

	if err := b.Status.IsValidStatus(); err != nil {
		return errors.New("invalid status")
	}
	return err

}

// func GetBucketCollection() *mongo.Collection {
// 	return config.Client.Database("storage").Collection("Bucket")
// }
