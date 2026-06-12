package redis

import (
	"context"
	"time"


	appusers "book_halal/internal/application/users/commands"

	"github.com/redis/go-redis/v9"
)


var _ appusers.VerificationRepository = (*VerificationRepository)(nil)

type VerificationRepository struct {
	client *redis.Client
}

func NewVerificationRepository(client *redis.Client) *VerificationRepository {
	return &VerificationRepository{client: client}
}

func (r *VerificationRepository) SaveCode(ctx context.Context, email string, code string) error {
	return r.client.Set(ctx, email, code, 5*time.Minute).Err()
}

func (r *VerificationRepository) GetCode(ctx context.Context, email string) (string, error) {
	
	return r.client.Get(ctx, email).Result()
}

func (r *VerificationRepository) DeleteCode(ctx context.Context, email string) error {
	return r.client.Del(ctx, email).Err()
}