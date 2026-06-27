package reader

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"book_halal/internal/domain/reader"
)

func (r *ReaderRepository) GetProgress(ctx context.Context, userID, bookID string) (*reader.ReadingProgress, error) {
	query := `
		SELECT user_id, book_id, page_number, updated_at
		FROM reading_progress
		WHERE user_id = $1 AND book_id = $2
	`

	var p reader.ReadingProgress
	err := r.pool.QueryRow(ctx, query, userID, bookID).Scan(
		&p.UserID, &p.BookID, &p.PageNumber, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // прогресс ещё не сохранялся
		}
		return nil, fmt.Errorf("get reading progress: %w", err)
	}

	return &p, nil
}
