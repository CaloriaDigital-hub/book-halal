package books

import (
	"context"
	"fmt"
	"time"

	book "book_halal/internal/domain/books/entity"
	repoBooks "book_halal/internal/domain/books"
)

func (r *BookRepository) GetPagesByBookID(ctx context.Context, bookID string) ([]book.Page, int, error) {
	query := `
		SELECT 
			bp.id, bp.book_id, bp.page_number, bp.image_url, bp.created_at, 
			b.total_pages
		FROM books b
		LEFT JOIN book_pages bp ON b.id = bp.book_id
		WHERE b.id = $1
		ORDER BY bp.page_number ASC
	`

	rows, err := r.pool.Query(ctx, query, bookID)
	if err != nil {
		return nil, 0, fmt.Errorf("query pages: %w", err)
	}
	defer rows.Close()

	var pages []book.Page = []book.Page{}
	totalPages := -1

	for rows.Next() {
		var (
			pageID     *string
			pBookID    *string
			pageNum    *int
			imageURL   *string
			createdAt  *time.Time 
			fetchedTotalPages int
		)

		err := rows.Scan(
			&pageID, &pBookID, &pageNum, &imageURL, &createdAt,
			&fetchedTotalPages,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan page: %w", err)
		}

		totalPages = fetchedTotalPages

		if pageID != nil {
			p := book.Page{
				ID:         *pageID,
				BookID:     *pBookID,
				PageNumber: *pageNum,
				ImageURL:   *imageURL,
			}
			
			
			if createdAt != nil {
				p.CreatedAt = *createdAt
			}	
			
			pages = append(pages, p)
		}
	}

	if totalPages == -1 {
		return nil, 0, fmt.Errorf("%w: id %s", repoBooks.ErrBookNotFound, bookID)
	}

	return pages, totalPages, nil
}