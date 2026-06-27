package valueobjects

import "errors"

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

var ErrInvalidRole = errors.New("invalid role")

func NewRole(val string) (Role, error) {
	switch Role(val) {
	case RoleUser, RoleAdmin:
		return Role(val), nil
	default:
		return "", ErrInvalidRole
	}
}

func (r Role) String() string {
	return string(r)
}

func (r Role) IsAdmin() bool {
	return r == RoleAdmin
}