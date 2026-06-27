package books

import (
	"context"
	"fmt"

	book "book_halal/internal/domain/books/entity"
)



func (r *BookRepository) GetPagesByBookIDPaginated(ctx context.Context, bookID string, offset, limit int) ([]book.Page, int, error) {
	var totalPages int
	err := r.pool.QueryRow(ctx,
		`SELECT total_pages FROM books WHERE id = $1`, bookID,
	).Scan(&totalPages)
	if err != nil {
		return nil, 0, fmt.Errorf("get total pages: %w", err)
	}

	query := `
		SELECT id, book_id, page_number, image_url, created_at
		FROM book_pages
		WHERE book_id = $1
		ORDER BY page_number ASC
		OFFSET $2
		LIMIT $3
	`

	rows, err := r.pool.Query(ctx, query, bookID, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("query paginated pages: %w", err)
	}
	defer rows.Close()

	pages := make([]book.Page, 0)
	for rows.Next() {
		var p book.Page
		if err := rows.Scan(&p.ID, &p.BookID, &p.PageNumber, &p.ImageURL, &p.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan page: %w", err)
		}
		pages = append(pages, p)
	}

	return pages, totalPages, nil
}
