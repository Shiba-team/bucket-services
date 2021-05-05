package model

import "errors"

type Permission string

func (p Permission) IsValidPermission() error {
	switch p {
	case "read", "write", "read:write":
		return nil
	}

	return errors.New("invalid permission")
}
