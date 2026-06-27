package books

import (
	"context"

	book "book_halal/internal/domain/books/entity"
)


func (r *BookRepository) Create(ctx context.Context, b *book.Book) error {
	query := `
		INSERT INTO books (id, title, author, description, price, total_pages, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
	`
	_, err := r.pool.Exec(ctx, query, b.ID, b.Title, b.Author, b.Description, b.Price.MinorValue(), b.TotalPages, b.Status)
	return err
}






