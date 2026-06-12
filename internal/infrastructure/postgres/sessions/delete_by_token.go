package sessions

import (
	"context"
	"fmt"
)

func (r *SessionRepository) DeleteByToken(ctx context.Context, token string) error {
	query := `DELETE FROM sessions WHERE token = $1`
	_, err := r.pool.Exec(ctx, query, token)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}