package reader

import (
	"context"
	"fmt"
)

func (r *ReaderRepository) UpsertProgress(ctx context.Context, userID, bookID string, page int) error {
	query := `
		INSERT INTO reading_progress (user_id, book_id, page_number, updated_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (user_id, book_id)
		DO UPDATE SET page_number = $3, updated_at = NOW()
	`

	_, err := r.pool.Exec(ctx, query, userID, bookID, page)
	if err != nil {
		return fmt.Errorf("upsert reading progress: %w", err)
	}

	return nil
}
