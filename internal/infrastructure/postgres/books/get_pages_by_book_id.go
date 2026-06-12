package books

import (
	"context"
	"fmt"

	book "book_halal/internal/domain/books/entity"
)

func (r *BookRepository) GetPagesByBookID(ctx context.Context, bookID string) ([]book.Page, int, error) {
	query := `
		SELECT bp.id, bp.book_id, bp.page_number, bp.image_url, bp.created_at, b.total_pages
		FROM book_pages bp
		JOIN books b ON b.id = bp.book_id
		WHERE bp.book_id = $1
		ORDER BY bp.page_number ASC
	`

	rows, err := r.pool.Query(ctx, query, bookID)
	if err != nil {
		return nil, 0, fmt.Errorf("query pages: %w", err)
	}
	defer rows.Close()

	var pages []book.Page
	var totalPages int

	for rows.Next() {
		var p book.Page
		err := rows.Scan(
			&p.ID, &p.BookID, &p.PageNumber, &p.ImageURL, &p.CreatedAt, &totalPages,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan page: %w", err)
		}
		pages = append(pages, p)
	}

	return pages, totalPages, nil
}