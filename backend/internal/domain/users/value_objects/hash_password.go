package valueobjects

import (
	"errors"

	pkgbcrypt "book_halal/internal/pkg/bcrypt"
)

var ErrInvalidPassword = errors.New("invalid password")

type HashedPassword struct {
	value string
}

func NewHashedPassword(password string) (HashedPassword, error) {
	if len(password) < 8 {
		return HashedPassword{}, errors.New("password must be at least 8 characters long")
	}
	if len(password) > 72 {
		return HashedPassword{}, errors.New("password is too long (max 72 bytes)")
	}

	hash, err := pkgbcrypt.Hash(password)
	if err != nil {
		return HashedPassword{}, err
	}

	return HashedPassword{value: hash}, nil
}

func NewHashedPasswordFromHash(hash string) HashedPassword {
	return HashedPassword{value: hash}
}

func (h HashedPassword) Compare(password string) error {
	if err := pkgbcrypt.Compare(h.value, password); err != nil {
		return ErrInvalidPassword
	}
	return nil
}

func (h HashedPassword) String() string {
	return h.value
}