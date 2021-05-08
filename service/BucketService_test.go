package service

import (
	"bucket/model"
	"testing"
)

func TestCreateBucket(t *testing.T) {

	bucket := model.NewFullBucket("testbucketname", "testuser", model.PermissionReadAndWrite, model.StatusPublic)

	CreateBucket(*bucket, func(err string) {
		t.Error("bucket name different testbucketname")
	}, func(id, name string) {
	})

}
