package reader

import (
	"context"
	"fmt"

	"book_halal/internal/domain/reader"
)

func (r *ReaderRepository) GetBookmarks(ctx context.Context, userID, bookID string) ([]reader.Bookmark, error) {
	query := `
		SELECT id, user_id, book_id, page_number, created_at
		FROM bookmarks
		WHERE user_id = $1 AND book_id = $2
		ORDER BY page_number ASC
	`

	rows, err := r.pool.Query(ctx, query, userID, bookID)
	if err != nil {
		return nil, fmt.Errorf("get bookmarks: %w", err)
	}
	defer rows.Close()

	var bookmarks []reader.Bookmark
	for rows.Next() {
		var b reader.Bookmark
		if err := rows.Scan(&b.ID, &b.UserID, &b.BookID, &b.PageNumber, &b.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan bookmark: %w", err)
		}
		bookmarks = append(bookmarks, b)
	}

	if bookmarks == nil {
		bookmarks = []reader.Bookmark{}
	}

	return bookmarks, nil
}
