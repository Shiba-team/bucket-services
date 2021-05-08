package model

import "testing"

func TestNewFile(t *testing.T) {
	file := NewFile("testfilename", "testfilelink", "testusername")

	if file.FileName != "testfilename" {
		t.Fail()
	}
}
