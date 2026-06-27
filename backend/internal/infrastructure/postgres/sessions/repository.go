package sessions

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"book_halal/internal/domain/sessions"
)

var _ sessions.Repository = (*SessionRepository)(nil)

type SessionRepository struct {
	pool *pgxpool.Pool
}

func NewSessionRepository(pool *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{pool: pool}
}