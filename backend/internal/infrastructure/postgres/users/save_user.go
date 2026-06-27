package users

import (
	"context"

	"book_halal/internal/domain/users/entity"
)

func (r *UserRepository) Save(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (id, first_name, last_name, email, password_hash)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE
		SET first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name,
			email = EXCLUDED.email,
			password_hash = EXCLUDED.password_hash


	`

	_, err := r.pool.Exec(ctx, query,
		user.ID.String(),
		user.FirstName,
		user.LastName,
		user.Email.String(),
		user.Password.String(),
	)

	return err
}