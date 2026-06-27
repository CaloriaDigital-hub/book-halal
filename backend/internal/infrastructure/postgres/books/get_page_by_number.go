package books

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	book "book_halal/internal/domain/books/entity"
)

// GetPageByNumber возвращает одну конкретную страницу книги по её номеру.
func (r *BookRepository) GetPageByNumber(ctx context.Context, bookID string, pageNumber int) (*book.Page, error) {
	query := `
		SELECT id, book_id, page_number, image_url, created_at
		FROM book_pages
		WHERE book_id = $1 AND page_number = $2
	`

	var p book.Page
	err := r.pool.QueryRow(ctx, query, bookID, pageNumber).Scan(
		&p.ID, &p.BookID, &p.PageNumber, &p.ImageURL, &p.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get page by number: %w", err)
	}

	return &p, nil
}
