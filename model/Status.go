package model

import "errors"

type Status string

func (p Status) IsValidStatus() error {
	switch p {
	case "public", "private":
		return nil
	}

	return errors.New("invalid status")
}
