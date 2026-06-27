package book

import (
	"context"

	"book_halal/internal/domain/books/entity"
)

type BookProcessor interface {
	Process(ctx context.Context, bookID string, pdfPath string) ([]entity.Page, error)
}