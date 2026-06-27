package valueobjects

import (
	uuid "book_halal/internal/pkg"
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidBookId = errors.New("invalid book id: cannot be empty")

type BookId struct {
	value string
}

func NewBookId(val string) (BookId, error) {
	val = strings.TrimSpace(val)

	if val == "" {
		return BookId{}, ErrInvalidBookId
	}

	return BookId{value: val}, nil
}

func (id BookId) String() string {
	return id.value
}

func (id BookId) IsEmpty() bool {
	return id.value == ""
}



func (id *BookId) Scan(value interface{}) error {
	if value == nil {
		return errors.New("book id cannot be null in database")
	}

	switch v := value.(type) {
	case string:
		id.value = v
		return nil
	case []byte:
		id.value = string(v)
		return nil
	case [16]byte: 
		
		id.value = uuid.FromBytes(v)
		return nil
	default:
		return fmt.Errorf("cannot scan type %T into BookId", value)
	}
}


func (id BookId) Value() (driver.Value, error) {
	if id.IsEmpty() {
		return nil, ErrInvalidBookId
	}
	return id.value, nil
}