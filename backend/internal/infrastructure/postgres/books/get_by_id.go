package books

import (
	"context"
	"fmt"

	"book_halal/internal/domain/books/entity"
	"book_halal/internal/domain/books/value_objects"
)

func (r *BookRepository) GetByID(ctx context.Context, bookID string) (*entity.Book, error) {
	query := `
		SELECT id, title, author, description, price, cover_url, total_pages, status
		FROM books
		WHERE id = $1 AND status = 'ready'
	`

	rows, err := r.pool.Query(ctx, query, bookID)
	if err != nil {
		return nil, fmt.Errorf("query book: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("book not found")
	}

	var b entity.Book
	var priceMinor int
	var description, coverURL *string

	err = rows.Scan(
		&b.ID, &b.Title, &b.Author, &description, &priceMinor, &coverURL,
		&b.TotalPages, &b.Status,
	)
	if err != nil {
		return nil, fmt.Errorf("scan book: %w", err)
	}

	if description != nil {
		b.Description = *description
	}
	if coverURL != nil {
		b.CoverURL = *coverURL
	}

	priceVO, _ := valueobjects.NewPriceFromMinor(priceMinor)
	b.Price = priceVO

	return &b, nil
}