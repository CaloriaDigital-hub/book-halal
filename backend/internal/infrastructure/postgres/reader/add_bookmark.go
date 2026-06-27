package reader

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	uuid "book_halal/internal/pkg"
	"book_halal/internal/domain/reader"
)

func (r *ReaderRepository) AddBookmark(ctx context.Context, userID, bookID string, page int) (*reader.Bookmark, error) {
	id := uuid.New()

	query := `
		INSERT INTO bookmarks (id, user_id, book_id, page_number, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT ON CONSTRAINT unique_bookmark DO NOTHING
		RETURNING id, user_id, book_id, page_number, created_at
	`

	var b reader.Bookmark
	err := r.pool.QueryRow(ctx, query, id, userID, bookID, page).Scan(
		&b.ID, &b.UserID, &b.BookID, &b.PageNumber, &b.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Закладка уже существует — вернём существующую
			return r.getExistingBookmark(ctx, userID, bookID, page)
		}
		return nil, fmt.Errorf("add bookmark: %w", err)
	}

	return &b, nil
}

func (r *ReaderRepository) getExistingBookmark(ctx context.Context, userID, bookID string, page int) (*reader.Bookmark, error) {
	query := `
		SELECT id, user_id, book_id, page_number, created_at
		FROM bookmarks
		WHERE user_id = $1 AND book_id = $2 AND page_number = $3
	`

	var b reader.Bookmark
	err := r.pool.QueryRow(ctx, query, userID, bookID, page).Scan(
		&b.ID, &b.UserID, &b.BookID, &b.PageNumber, &b.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get existing bookmark: %w", err)
	}

	return &b, nil
}
