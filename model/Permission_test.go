package model

import (
	"testing"
)

func TestIsValidPermission(t *testing.T) {
	var p Permission
	p = PermissionRead

	p.IsValidPermission(func(s string) {
		t.Fail()
	})

	p = 10

	p.IsValidPermission(func(s string) {
		t.Fail()
	})

}
