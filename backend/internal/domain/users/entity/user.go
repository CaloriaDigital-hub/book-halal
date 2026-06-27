package entity

import (
	"book_halal/internal/domain/users/value_objects"
)

type User struct {
	ID       valueobjects.UserId
	FirstName string
	LastName  string
	Password  valueobjects.HashedPassword
	Email     valueobjects.Email
	Role      valueobjects.Role
}

func NewUser(id valueobjects.UserId, firstName, lastName string, email valueobjects.Email, passHash valueobjects.HashedPassword) *User {
	return &User{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  passHash,
		Role:      valueobjects.RoleUser,
	}
}

// Reconstruct rebuilds a User from persisted data (no defaults applied).
func Reconstruct(id valueobjects.UserId, firstName, lastName string, email valueobjects.Email, passHash valueobjects.HashedPassword, role valueobjects.Role) *User {
	return &User{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  passHash,
		Role:      role,
	}
}