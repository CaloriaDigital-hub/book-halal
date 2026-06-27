package commands

import (
	"context"
	"errors"
)

// ErrCodeNotFound is returned by VerificationRepository.GetCode when the
// OTP key does not exist in the store (never saved or TTL expired).
var ErrCodeNotFound = errors.New("verification code not found or expired")

// VerificationRepository — порт для хранения OTP-кодов верификации (Redis).
type VerificationRepository interface {
	SaveCode(ctx context.Context, email string, code string) error
	GetCode(ctx context.Context, email string) (string, error)
	DeleteCode(ctx context.Context, email string) error
}

// EmailSender — порт для отправки email-уведомлений.
type EmailSender interface {
	SendOTP(ctx context.Context, toEmail string, code string) error
}
