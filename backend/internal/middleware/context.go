package middleware

import (
	"context"

	"book_halal/internal/domain/users/entity"
)

type contextKey string

const userContextKey contextKey = "user"

func WithUser(ctx context.Context, user *entity.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func UserFromContext(ctx context.Context) (*entity.User, bool) {
	user, ok := ctx.Value(userContextKey).(*entity.User)
	return user, ok
}