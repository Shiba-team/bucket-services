package model

import "time"

type File struct {
	Id         string    `json:"_id"`
	FileLink   string    `json:"filelink"`
	CreatedBy  string    `json:"createdby"`
	FileName   string    `json:"filename"`
	CreatedAt  time.Time `json:"createdat"`
	LastModify time.Time `json:"lastmodify"`
	Status     string    `json:"status"`
	Permission string    `json:"permission"`
}

func NewFile(filename string, filelink string, username string) File {
	var f File
	f.FileName = filename
	f.FileLink = filelink
	f.CreatedAt = time.Now()
	f.CreatedBy = username
	f.LastModify = time.Now()
	f.Status = "private"
	f.Permission = "read:write"
	return f
}
