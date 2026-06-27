package users

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"book_halal/internal/domain/users/entity"
	"book_halal/internal/domain/users/value_objects"
)

func (r *UserRepository) FindByEmail(ctx context.Context, email valueobjects.Email) (*entity.User, error) {
	query := `
		SELECT id, first_name, last_name, password_hash, role
		FROM users
		WHERE email = $1
		LIMIT 1
	`

	var id, firstName, lastName, passHash, roleStr string

	err := r.pool.QueryRow(ctx, query, email.String()).Scan(&id, &firstName, &lastName, &passHash, &roleStr)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	userID, _ := valueobjects.NewUserId(id)
	password := valueobjects.NewHashedPasswordFromHash(passHash)
	role, err := valueobjects.NewRole(roleStr)
	if err != nil {
		return nil, err
	}

	return entity.Reconstruct(userID, firstName, lastName, email, password, role), nil
}