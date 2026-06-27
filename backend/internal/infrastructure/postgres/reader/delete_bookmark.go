package reader

import (
	"context"
	"fmt"
)

func (r *ReaderRepository) DeleteBookmark(ctx context.Context, userID, bookID string, page int) error {
	query := `
		DELETE FROM bookmarks
		WHERE user_id = $1 AND book_id = $2 AND page_number = $3
	`

	tag, err := r.pool.Exec(ctx, query, userID, bookID, page)
	if err != nil {
		return fmt.Errorf("delete bookmark: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("bookmark not found")
	}

	return nil
}
