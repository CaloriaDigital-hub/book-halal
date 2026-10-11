package purchases

import (
	"errors"
	"testing"
)

func TestNewContact(t *testing.T) {
	cases := []struct {
		typ, value, want string
		ok               bool
	}{
		{"phone", "+996 555 123 456", "+996 555 123 456", true},
		{"whatsapp", "0555123456", "0555123456", true},
		{"telegram", "durov_book", "@durov_book", true},
		{"telegram", "@durov_book", "@durov_book", true},
		{"telegram", "+996555123456", "+996555123456", true},
		{"phone", "abc", "", false},
		{"telegram", "@ab", "", false},
		{"email", "a@b.c", "", false},
		{"phone", "", "", false},
	}
	for _, c := range cases {
		got, err := NewContact(c.typ, c.value)
		if c.ok {
			if err != nil || got.Value != c.want || string(got.Type) != c.typ {
				t.Errorf("NewContact(%q, %q) = %+v, %v; want %q", c.typ, c.value, got, err, c.want)
			}
			continue
		}
		if !errors.Is(err, ErrInvalidContact) {
			t.Errorf("NewContact(%q, %q) err = %v; want ErrInvalidContact", c.typ, c.value, err)
		}
	}
}
