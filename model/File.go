package model

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

func RandomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

type File struct {
	FileLink     string     `json:"filelink"  bson:"filelink"`
	DownloadLink string     `json:"downloadlink"  bson:"downloadlink"`
	CreatedBy    string     `json:"createdby" bson:"createdby"`
	FileName     string     `json:"filename" bson:"filename"`
	S3Name       string     `json:"s3name" bson:"s3name"`
	CreatedAt    time.Time  `json:"createdat" bson:"createdat"`
	LastModify   time.Time  `json:"lastmodify" bson:"lastmodify"`
	Status       Status     `json:"status" bson:"status"`
	Permission   Permission `json:"permission" bson:"permission"`
	Size         int64      `json:"size" bson:"size"`
}

func NewFile(filename string, filelink string, username string) File {
	var f File
	f.FileName = filename
	f.FileLink = filelink
	f.CreatedAt = time.Now()
	f.CreatedBy = username
	f.LastModify = time.Now()
	f.Status = StatusPrivate
	f.Permission = PermissionRead
	f.Size = 0
	return f
}
