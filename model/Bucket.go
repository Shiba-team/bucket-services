package model

import "time"

type Bucket struct {
	Id           string    `json:"_id"`
	Owner        string    `json:"owner"`
	BucketName   string    `json:"bucketname" validate:"required" binding:"required"`
	CreatedAt    time.Time `json:"createdat"`
	ListFile     []File    `json:"listfile"`
	Status       string    `json:"status"`
	Permission   string    `json:"permission"`
	LastModified time.Time `json:"lastmodified"`
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

// func GetBucketCollection() *mongo.Collection {
// 	return config.Client.Database("storage").Collection("Bucket")
// }
