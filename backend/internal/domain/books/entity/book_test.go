package entity_test

import (
	"testing"

	"book_halal/internal/domain/books/entity"
	valueobjects "book_halal/internal/domain/books/value_objects"
)

func TestNew(t *testing.T) {
	// Собираем честный VO через конструктор
	validID, err := valueobjects.NewBookId("book-123")
	if err != nil {
		t.Fatalf("не удалось создать VO BookId для теста: %v", err)
	}

	validPrice, err := valueobjects.NewPriceFromMinor(100)
	
	if err != nil {
		t.Fatalf("не удалось создать VO Price для теста: %v", err)
	}// Если Price тоже станет VO — переделать аналогично

	tests := []struct {
		name        string
		id          valueobjects.BookId
		title       string
		author      string
		description string
		price       valueobjects.Price
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "успешное создание книги",
			id:          validID,
			title:       "Чистая Архитектура",
			author:      "Роберт Мартин",
			description: "Искусство разработки ПО",
			price:       validPrice,
			wantErr:     false,
		},
		{
			name:        "ошибка: пустой title",
			id:          validID,
			title:       "",
			author:      "Роберт Мартин",
			description: "Описание",
			price:       validPrice,
			wantErr:     true,
			errMsg:      "book title cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			book, err := entity.New(tt.id, tt.title, tt.author, tt.description, tt.price)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("ожидали ошибку, но получили nil")
				}
				if err.Error() != tt.errMsg {
					t.Errorf("got error msg %q, want %q", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Fatalf("не ожидали ошибку, но получили: %v", err)
			}

			// Сравниваем VO
			if book.ID != tt.id {
				t.Errorf("got ID %v, want %v", book.ID, tt.id)
			}
		})
	}
}