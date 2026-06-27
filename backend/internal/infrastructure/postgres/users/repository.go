package users

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"book_halal/internal/domain/users"
)


var _ users.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	pool *pgxpool.Pool
}


func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}