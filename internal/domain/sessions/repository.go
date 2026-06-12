package sessions

import (
	"context"
	"book_halal/internal/domain/sessions/entity"
)

type Repository interface {
	Save(ctx context.Context, session *entity.Session) error
	FindByToken(ctx context.Context, token string) (*entity.Session, error)
	DeleteByToken(ctx context.Context, token string) error
}