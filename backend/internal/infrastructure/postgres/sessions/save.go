package sessions

import (
	"context"
	"fmt"

	"book_halal/internal/domain/sessions/entity"
)

func (r *SessionRepository) Save(ctx context.Context, s *entity.Session) error {
	query := `
		INSERT INTO sessions (id, user_id, token, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(ctx, query, s.ID, s.UserID, s.Token, s.ExpiresAt, s.CreatedAt)
	if err != nil {
		return fmt.Errorf("save session: %w", err)
	}
	return nil
}