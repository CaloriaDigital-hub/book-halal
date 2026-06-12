package valueobjects


import (
	"errors"
	"strings"
)


var ErrInvalidEmail = errors.New("invalid email")

type Email struct {
	value string

}

func NewEmail(val string) (Email, error) {
	val = strings.TrimSpace(strings.ToLower(val))

	if !strings.Contains(val, "@") {
		return Email{}, ErrInvalidEmail
	}

	return Email{value: val}, nil

}

func (e Email) String() string {
	return e.value
}