package valueobjects_test

import (
	"testing"

	valueobjects "book_halal/internal/domain/books/value_objects"
)

func TestNewPriceFromMajor(t *testing.T) {
	tests := []struct {
		name        string
		majorAmount int
		wantMinor   int
		wantErr     bool
	}{
		{
			name:        "валидная позитивная цена",
			majorAmount: 15,
			wantMinor:   1500,
			wantErr:     false,
		},
		{
			name:        "бесплатный товар (0)",
			majorAmount: 0,
			wantMinor:   0,
			wantErr:     false,
		},
		{
			name:        "ошибка: отрицательная цена",
			majorAmount: -10,
			wantMinor:   0,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			price, err := valueobjects.NewPriceFromMajor(tt.majorAmount)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("ожидали ошибку для majorAmount = %d, получили nil", tt.majorAmount)
				}
				return
			}

			if err != nil {
				t.Fatalf("не ожидали ошибку, но получили: %v", err)
			}

			if price.MinorValue() != tt.wantMinor {
				t.Errorf("got MinorValue() = %d, want %d", price.MinorValue(), tt.wantMinor)
			}
		})
	}
}

func TestNewPriceFromMinor(t *testing.T) {
	tests := []struct {
		name        string
		minorAmount int
		wantMajor   int
		wantErr     bool
	}{
		{
			name:        "валидная цена в копейках",
			minorAmount: 2500,
			wantMajor:   25,
			wantErr:     false,
		},
		{
			name:        "цена с копейками (дробная часть отбрасывается в MajorValue)",
			minorAmount: 2599, // 25.99
			wantMajor:   25,   // 2599 / 100 = 25
			wantErr:     false,
		},
		{
			name:        "ошибка: отрицательные копейки",
			minorAmount: -500,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			price, err := valueobjects.NewPriceFromMinor(tt.minorAmount)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("ожидали ошибку для minorAmount = %d, получили nil", tt.minorAmount)
				}
				return
			}

			if err != nil {
				t.Fatalf("не ожидали ошибку: %v", err)
			}

			if price.MajorValue() != tt.wantMajor {
				t.Errorf("got MajorValue() = %d, want %d", price.MajorValue(), tt.wantMajor)
			}
		})
	}
}

func TestPrice_IsFree(t *testing.T) {
	freePrice, err := valueobjects.NewPriceFromMajor(0)
	if err != nil {
		t.Fatalf("failed to create free price: %v", err)
	}

	paidPrice, err := valueobjects.NewPriceFromMajor(10)
	if err != nil {
		t.Fatalf("failed to create paid price: %v", err)
	}

	if !freePrice.IsFree() {
		t.Errorf("expected price with 0 amount to be free")
	}

	if paidPrice.IsFree() {
		t.Errorf("expected price with 1000 minor amount NOT to be free")
	}
}