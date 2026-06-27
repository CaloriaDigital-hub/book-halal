package commands

import (
	"context"
	"errors"
	"fmt"
	"time"

	sessionDomain "book_halal/internal/domain/sessions"
	sessionEntity "book_halal/internal/domain/sessions/entity"
	userDomain "book_halal/internal/domain/users"
	"book_halal/internal/domain/users/value_objects"
	pkgtoken "book_halal/internal/pkg/token"
	pkguuid "book_halal/internal/pkg"
)

// --- DTO ---

type SignInCommand struct {
	Email    string
	Password string
}

type SignInResult struct {
	Token     string
	ExpiresAt string
}

// --- Interface ---

type SignInHandler interface {
	Handle(ctx context.Context, cmd SignInCommand) (*SignInResult, error)
}

// --- Handler ---

var ErrInvalidCredentials = errors.New("invalid email or password")

const sessionTTL = 30 * 24 * time.Hour

type signInHandler struct {
	userRepo    userDomain.UserRepository
	sessionRepo sessionDomain.Repository
}

func NewSignInHandler(userRepo userDomain.UserRepository, sessionRepo sessionDomain.Repository) SignInHandler {
	return &signInHandler{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (h *signInHandler) Handle(ctx context.Context, cmd SignInCommand) (*SignInResult, error) {
	emailVO, err := valueobjects.NewEmail(cmd.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	user, err := h.userRepo.FindByEmail(ctx, emailVO)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := user.Password.Compare(cmd.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := pkgtoken.New()
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	expiresAt := time.Now().Add(sessionTTL)
	session := sessionEntity.NewSession(pkguuid.New(), user.ID.String(), token, expiresAt)

	if err := h.sessionRepo.Save(ctx, session); err != nil {
		return nil, fmt.Errorf("save session: %w", err)
	}

	return &SignInResult{
		Token:     token,
		ExpiresAt: expiresAt.Format(time.RFC3339),
	}, nil
}
