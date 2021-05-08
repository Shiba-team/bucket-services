package model

import (
	"testing"
)

func TestNewFullBucket(t *testing.T) {
	testFullBucketSuccess := NewFullBucket("testname", "testuser", PermissionRead, StatusPublic)

	if testFullBucketSuccess.BucketName != "testname" ||
		testFullBucketSuccess.Owner != "testuser" ||
		testFullBucketSuccess.Permission != PermissionRead ||
		testFullBucketSuccess.Status != StatusPublic {
		t.Fail()
	}
}

func TestIsValidBucket(t *testing.T) {
	testbucket := NewBucket("testbucket", "testuser")
	testbucket.IsValidBucket(func(s string) {
		{
			t.Fail()
		}
	}, func() {})

	testFullBucketSuccess := NewFullBucket("testname", "testuser", PermissionRead, StatusPublic)
	testFullBucketSuccess.IsValidBucket(func(s string) {
		{
			t.Fail()
		}
	}, func() {})

	testFullBucketFailedPermission := NewFullBucket("testname", "testuser", 10, StatusPublic)
	testFullBucketFailedPermission.IsValidBucket(func(s string) {
		{
			t.Fail()
		}
	}, func() {})

	testFullBucketFailedStatus := NewFullBucket("testname", "testuser", PermissionRead, 10)
	testFullBucketFailedStatus.IsValidBucket(func(s string) {
		{
			t.Fail()
		}
	}, func() {})
}

func TestNewBucket(t *testing.T) {
	testbucket := NewBucket("testbucket", "testuser")
	testbucket.IsValidBucket(func(s string) {
		{
			t.Fail()
		}
	}, func() {})

	if testbucket.BucketName != "testbucket" || testbucket.Owner != "testuser" {
		t.Fail()
	}
}
