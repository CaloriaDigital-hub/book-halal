package valueobjects

import (
	"errors"
	"strings"
)

var ErrInvalidIserId = errors.New("invalid user id: cannot be empty")

type UserId struct {
	value string
}


func NewUserId(val string) (UserId, error) {
	val = strings.TrimSpace(val)

	if val == "" {
		return UserId{}, ErrInvalidIserId
	}


	
	return UserId{value: val}, nil
}

func (id UserId) String() string {
	return id.value
}

func (id UserId) IsEmpty() bool {
	return id.value == ""
}