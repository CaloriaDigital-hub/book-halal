package books


import (
	"context"

	book "book_halal/internal/domain/books/entity"
)

func (r *BookRepository) UpdateStatusWithPageCount(ctx context.Context, bookID string, status book.Status, totalPages int) error {
	query := `UPDATE books SET status = $1, total_pages = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.pool.Exec(ctx, query, status, totalPages, bookID)
	return err
}