package redis

import (
	"context"
	"errors"
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

// GetCode returns the OTP code for the given email.
// Returns appusers.ErrCodeNotFound if the key does not exist or has expired.
func (r *VerificationRepository) GetCode(ctx context.Context, email string) (string, error) {
	val, err := r.client.Get(ctx, email).Result()
	if errors.Is(err, redis.Nil) {
		// Key doesn't exist or TTL expired — not a Redis error, just "not found".
		return "", appusers.ErrCodeNotFound
	}
	return val, err
}

func (r *VerificationRepository) DeleteCode(ctx context.Context, email string) error {
	return r.client.Del(ctx, email).Err()
}