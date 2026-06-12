package handlers

import (
	uuid "book_halal/internal/pkg"
	"book_halal/internal/application/users/commands"
	"book_halal/internal/domain/users"
	"book_halal/internal/domain/users/entity"
	valueobjects "book_halal/internal/domain/users/value_objects"
	"context"
	"errors"
)

type CreateUserHandler struct {
	userRepo users.UserRepository
}

func NewUserCreateHandler(repo users.UserRepository) *CreateUserHandler {
	return &CreateUserHandler{userRepo: repo}
}

func (h *CreateUserHandler) Handle(ctx context.Context, cmd commands.CreateUserCommand) error {

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
