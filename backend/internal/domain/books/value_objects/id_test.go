package valueobjects_test

import (
	"testing"

	valueobjects "book_halal/internal/domain/books/value_objects"
)

func TestNewBookId(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "валидный id",
			input:   "book-123",
			want:    "book-123",
			wantErr: false,
		},
		{
			name:    "id с пробелами (должен сработать TrimSpace)",
			input:   "   book-456   ",
			want:    "book-456",
			wantErr: false,
		},
		{
			name:    "ошибка: пустая строка",
			input:   "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "ошибка: строки из одних пробелов",
			input:   "     ",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := valueobjects.NewBookId(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("ожидали ошибку, но получили nil")
				}
				if err != valueobjects.ErrInvalidBookId {
					t.Errorf("got err %v, want %v", err, valueobjects.ErrInvalidBookId)
				}
				return
			}

			if err != nil {
				t.Fatalf("не ожидали ошибку, но получили: %v", err)
			}

			if id.String() != tt.want {
				t.Errorf("got String() = %q, want %q", id.String(), tt.want)
			}

			if id.IsEmpty() {
				t.Errorf("IsEmpty() вернул true для валидного ID")
			}
		})
	}
}

func TestBookId_Scan(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    string
		wantErr bool
	}{
		{
			name:    "scan из string",
			input:   "uuid-string-123",
			want:    "uuid-string-123",
			wantErr: false,
		},
		{
			name:    "scan из []byte",
			input:   []byte("uuid-bytes-123"),
			want:    "uuid-bytes-123",
			wantErr: false,
		},
		{
			name:    "ошибка: nil значение",
			input:   nil,
			wantErr: true,
		},
		{
			name:    "ошибка: неподдерживаемый тип (int)",
			input:   12345,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var id valueobjects.BookId
			err := id.Scan(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("ожидали ошибку при Scan, но получили nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("не ожидали ошибку при Scan: %v", err)
			}

			if id.String() != tt.want {
				t.Errorf("got %q, want %q", id.String(), tt.want)
			}
		})
	}
}

func TestBookId_Value(t *testing.T) {
	t.Run("успешное получение driver.Value", func(t *testing.T) {
		id, err := valueobjects.NewBookId("valid-id")
		if err != nil {
			t.Fatalf("failed to create id: %v", err)
		}

		val, err := id.Value()
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		if val != "valid-id" {
			t.Errorf("got driver.Value = %v, want %v", val, "valid-id")
		}
	})

	t.Run("ошибка при сухом/пустом BookId", func(t *testing.T) {
		var emptyID valueobjects.BookId

		val, err := emptyID.Value()
		if err == nil {
			t.Fatalf("ожидали ошибку для пустого BookId, получили val = %v", val)
		}
		if err != valueobjects.ErrInvalidBookId {
			t.Errorf("got err %v, want %v", err, valueobjects.ErrInvalidBookId)
		}
	})
}