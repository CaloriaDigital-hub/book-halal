package handlers

import (
	"context"
	"errors"

	uuid "book_halal/internal/pkg"
	"book_halal/internal/application/users/commands"
	"book_halal/internal/domain/users/entity"
	valueobjects "book_halal/internal/domain/users/value_objects"

	appusers "book_halal/internal/application/users/commands"
	domainusers "book_halal/internal/domain/users"
)

type ConfirmRegistrationHandler struct {
	userRepo domainusers.UserRepository
	codeRepo appusers.VerificationRepository
}

func NewConfirmRegistrationHandler(
	userRepo domainusers.UserRepository,
	codeRepo appusers.VerificationRepository,
) *ConfirmRegistrationHandler {
	return &ConfirmRegistrationHandler{
		userRepo: userRepo,
		codeRepo: codeRepo,
	}
}

func (h *ConfirmRegistrationHandler) Handle(ctx context.Context, cmd commands.ConfirmRegistrationCommand) error {
	// 1. Достаем код из Redis
	savedCode, err := h.codeRepo.GetCode(ctx, cmd.Email)
	if err != nil {
		return errors.New("verification code expired or not found")
	}

	// 2. Проверяем совпадение
	if savedCode != cmd.Code {
		return errors.New("invalid verification code")
	}

	// 3. Формируем доменные объекты
	userID, err := valueobjects.NewUserId(uuid.New())
	if err != nil {
		return err
	}

	email, err := valueobjects.NewEmail(cmd.Email)
	if err != nil {
		return err
	}

	password, err := valueobjects.NewHashedPassword(cmd.Password)
	if err != nil {
		return err
	}

	// 4. Создаем чистую сущность
	user := entity.NewUser(userID, cmd.FirstName, cmd.LastName, email, password)

	// 5. Сохраняем в PostgreSQL
	if err := h.userRepo.Save(ctx, user); err != nil {
		return err
	}

	// 6. Сжигаем код в Redis, чтобы его нельзя было использовать дважды
	_ = h.codeRepo.DeleteCode(ctx, cmd.Email)

	return nil
}
