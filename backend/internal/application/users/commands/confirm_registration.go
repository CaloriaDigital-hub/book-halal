package commands

import (
	"context"
	"errors"
	"fmt"

	uuid "book_halal/internal/pkg"
	"book_halal/internal/domain/users/entity"
	valueobjects "book_halal/internal/domain/users/value_objects"
	domainusers "book_halal/internal/domain/users"
)

// --- DTO ---

type ConfirmRegistrationCommand struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	Code      string
}

// --- Interface ---

type ConfirmRegistrationHandler interface {
	Handle(ctx context.Context, cmd ConfirmRegistrationCommand) error
}

// --- Handler ---

type confirmRegistrationHandler struct {
	userRepo domainusers.UserRepository
	codeRepo VerificationRepository
}

func NewConfirmRegistrationHandler(
	userRepo domainusers.UserRepository,
	codeRepo VerificationRepository,
) ConfirmRegistrationHandler {
	return &confirmRegistrationHandler{
		userRepo: userRepo,
		codeRepo: codeRepo,
	}
}

func (h *confirmRegistrationHandler) Handle(ctx context.Context, cmd ConfirmRegistrationCommand) error {
	// 1. Достаем код из Redis
	savedCode, err := h.codeRepo.GetCode(ctx, cmd.Email)
	if err != nil {
		if errors.Is(err, ErrCodeNotFound) {
			// Нормальная ситуация: код истёк или не был запрошен
			return errors.New("verification code expired or not found")
		}
		// Реальная ошибка Redis (недоступен, timeout и т.д.)
		return fmt.Errorf("failed to get verification code: %w", err)
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
