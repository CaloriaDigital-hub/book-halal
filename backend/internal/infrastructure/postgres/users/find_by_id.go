package users

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"book_halal/internal/domain/users/entity"
	"book_halal/internal/domain/users/value_objects"
)

func (r *UserRepository) FindByID(ctx context.Context, id valueobjects.UserId) (*entity.User, error) {
	query := `
		SELECT id, first_name, last_name, email, password_hash, role
		FROM users
		WHERE id = $1
		LIMIT 1
	`

	var idStr, firstName, lastName, emailStr, passHash, roleStr string

	err := r.pool.QueryRow(ctx, query, id.String()).Scan(&idStr, &firstName, &lastName, &emailStr, &passHash, &roleStr)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	userID, _ := valueobjects.NewUserId(idStr)
	email, _ := valueobjects.NewEmail(emailStr)
	password := valueobjects.NewHashedPasswordFromHash(passHash)
	role, err := valueobjects.NewRole(roleStr)
	if err != nil {
		return nil, err
	}

	return entity.Reconstruct(userID, firstName, lastName, email, password, role), nil
}