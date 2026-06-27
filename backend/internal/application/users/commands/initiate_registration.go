package commands

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"book_halal/internal/domain/users/value_objects"
	domainusers "book_halal/internal/domain/users"
)

// --- DTO ---

type InitiateRegistrationCommand struct {
	Email string
}

// --- Interface ---

type InitiateRegistrationHandler interface {
	Handle(ctx context.Context, cmd InitiateRegistrationCommand) error
}

// --- Handler ---

type initiateRegistrationHandler struct {
	userRepo domainusers.UserRepository
	codeRepo VerificationRepository
	emailSvc EmailSender
}

func NewInitiateRegistrationHandler(
	userRepo domainusers.UserRepository,
	codeRepo VerificationRepository,
	emailSvc EmailSender,
) InitiateRegistrationHandler {
	return &initiateRegistrationHandler{
		userRepo: userRepo,
		codeRepo: codeRepo,
		emailSvc: emailSvc,
	}
}

func (h *initiateRegistrationHandler) Handle(ctx context.Context, cmd InitiateRegistrationCommand) error {

	email, err := valueobjects.NewEmail(cmd.Email)
	if err != nil {
		return err
	}

	existingUser, err := h.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return fmt.Errorf("user with this email already exists")
	}

	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return err
	}
	code := fmt.Sprintf("%06d", n.Int64()+100000)

	if err := h.codeRepo.SaveCode(ctx, email.String(), code); err != nil {
		return err
	}

	if err := h.emailSvc.SendOTP(ctx, email.String(), code); err != nil {
		return err
	}

	return nil
}
