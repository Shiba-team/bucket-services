package service

import (
	"bucket/model"
	"testing"
)

func TestCreateBucket(t *testing.T) {

	bucket := model.NewFullBucket("testbucketname", "testuser", "read", "private")

	result, err := CreateBucket(*bucket)

	if err == nil && result != "testbucketname" {
		t.Error("bucket name different testbucketname")
	}

}
