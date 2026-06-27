package users

import (
	"context"
	"book_halal/internal/domain/users/entity"
	"book_halal/internal/domain/users/value_objects"
)

type UserRepository interface {
	Save(ctx context.Context, user *entity.User) error
	FindByEmail(ctx context.Context, email valueobjects.Email) (*entity.User, error)
	FindByID(ctx context.Context, id valueobjects.UserId) (*entity.User, error)
	

}
