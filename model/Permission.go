package model

type Permission int

const (
	PermissionRead Permission = iota
	PermissionWrite
	PermissionReadAndWrite
)

func (p Permission) IsValidPermission(iferror func(string)) bool {
	switch p {
	case PermissionRead, PermissionWrite, PermissionReadAndWrite:
		return true
	}

	iferror("invalid permission")
	return false
}
