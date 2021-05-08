package model

type Status int

const (
	StatusPublic Status = iota
	StatusPrivate
)

func (p Status) IsValidStatus(iferror func(err string)) bool {
	switch p {
	case StatusPublic, StatusPrivate:
		return true
	}

	iferror("invalid Status Code")
	return false
}
