package commands

import (
	"context"

)


type VerificationRepository interface {
	SaveCode(ctx context.Context, email string, code string) error
	GetCode(ctx context.Context, email string) (string, error)
	DeleteCode(ctx context.Context, email string) error
}


type EmailSender interface {
	SendOTP(ctx context.Context, toEmail string, code string) error
}

type SignInHandler interface {
	Handle(ctx context.Context, cmd SignInCommand) (*SignInResult, error)
}