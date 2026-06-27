package uuid

import "github.com/google/uuid"


func New() string {
	return uuid.New().String()
}

func FromBytes(b [16]byte) string {
    return uuid.UUID(b).String()
}

func FromString(s string) (uuid.UUID, error) {
    return uuid.Parse(s)
}

func MustParse(s string) uuid.UUID {
    return uuid.MustParse(s)
}