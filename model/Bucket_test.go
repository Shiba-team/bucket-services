package model

import (
	"testing"
)

func TestNewFullBucket(t *testing.T) {
	testFullBucketSuccess := NewFullBucket("testname", "testuser", "read", "public")

	if testFullBucketSuccess.BucketName != "testname" ||
		testFullBucketSuccess.Owner != "testuser" ||
		testFullBucketSuccess.Permission != "read" ||
		testFullBucketSuccess.Status != "public" {
		t.Fail()
	}
}

func TestIsValidBucket(t *testing.T) {
	testbucket := NewBucket("testbucket", "testuser")
	if err := testbucket.IsValidBucket(); err != nil {
		t.Fail()
	}

	testFullBucketSuccess := NewFullBucket("testname", "testuser", "read", "public")

	if err := testFullBucketSuccess.IsValidBucket(); err != nil {
		t.Error(err)
	}

	testFullBucketFailedPermission := NewFullBucket("testname", "testuser", "readad", "public")

	if err := testFullBucketFailedPermission.IsValidBucket(); err == nil || err.Error() != "invalid permission" {
		t.Error(err)
	}

	testFullBucketFailedStatus := NewFullBucket("testname", "testuser", "read", "pri")

	if err := testFullBucketFailedStatus.IsValidBucket(); err == nil || err.Error() != "invalid status" {
		t.Error(err)
	}
}

func TestNewBucket(t *testing.T) {
	testbucket := NewBucket("testbucket", "testuser")
	if err := testbucket.IsValidBucket(); err != nil {
		t.Fail()
	}

	if testbucket.BucketName != "testbucket" || testbucket.Owner != "testuser" {
		t.Fail()
	}
}
