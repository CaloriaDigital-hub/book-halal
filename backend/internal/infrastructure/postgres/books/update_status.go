package books

import (
	"context"

	book "book_halal/internal/domain/books/entity"
)


func (r *BookRepository) UpdateStatus(ctx context.Context, bookID string, status book.Status) error {
	query := `UPDATE books SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, status, bookID)
	return err
}