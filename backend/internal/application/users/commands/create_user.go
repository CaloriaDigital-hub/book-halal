package commands

import (
	"context"
	"errors"

	uuid "book_halal/internal/pkg"
	"book_halal/internal/domain/users"
	"book_halal/internal/domain/users/entity"
	valueobjects "book_halal/internal/domain/users/value_objects"
)

// --- DTO ---

type CreateUserCommand struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

// --- Interface ---

type CreateUserHandler interface {
	Handle(ctx context.Context, cmd CreateUserCommand) error
}

// --- Handler ---

type createUserHandler struct {
	userRepo users.UserRepository
}

func NewCreateUserHandler(repo users.UserRepository) CreateUserHandler {
	return &createUserHandler{userRepo: repo}
}

func (h *createUserHandler) Handle(ctx context.Context, cmd CreateUserCommand) error {

	userID, err := valueobjects.NewUserId(uuid.New())
	if err != nil {
		return err
	}

	email, err := valueobjects.NewEmail(cmd.Email)
	if err != nil {
		return err
	}

	passHash, err := valueobjects.NewHashedPassword(cmd.Password)
	if err != nil {
		return err
	}

	if existingUser, _ := h.userRepo.FindByEmail(ctx, email); existingUser != nil {
		return errors.New("user with this email already exists")
	}

	user := entity.NewUser(userID, cmd.FirstName, cmd.LastName, email, passHash)

	return h.userRepo.Save(ctx, user)
}
