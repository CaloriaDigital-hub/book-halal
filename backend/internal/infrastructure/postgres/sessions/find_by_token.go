package sessions

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"book_halal/internal/domain/sessions/entity"
)

func (r *SessionRepository) FindByToken(ctx context.Context, token string) (*entity.Session, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at
		FROM sessions
		WHERE token = $1
	`

	var s entity.Session
	err := r.pool.QueryRow(ctx, query, token).Scan(
		&s.ID, &s.UserID, &s.Token, &s.ExpiresAt, &s.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("find session: %w", err)
	}
	return &s, nil
}