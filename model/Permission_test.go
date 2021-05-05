package model

import (
	"testing"
)

func TestIsValidPermission(t *testing.T) {
	var p Permission
	p = "read"

	if err := p.IsValidPermission(); err != nil {
		t.Fail()
	}

	p = "test"

	if err := p.IsValidPermission(); err == nil {
		t.Fail()
	}

}
