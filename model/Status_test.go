package model

import (
	"fmt"
	"testing"
)

func TestIsValidStatus(t *testing.T) {
	var p Status
	p = StatusPublic

	p.IsValidStatus(func(err string) {
		fmt.Print(err)
		t.Error(err)
	})

	p = 2

	p.IsValidStatus(func(err string) {
		fmt.Print(err)
		t.Error(err)
	})

}
