package model

import (
	"fmt"
	"testing"
)

func TestIsValidStatus(t *testing.T) {
	var p Status
	p = "public"

	if err := p.IsValidStatus(); err != nil {
		fmt.Print(err)
		t.Error(err)
	}

	p = "test"

	if err := p.IsValidStatus(); err == nil {
		fmt.Print(err)
		t.Error(err)
	}

}
