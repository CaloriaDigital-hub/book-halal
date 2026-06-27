package controllers

import (
	"context"

	"book_halal/internal/application/users/commands"
)

type InitiateRegistrationUseCase interface {
	Handle(ctx context.Context, cmd commands.InitiateRegistrationCommand) error
}

type ConfirmRegistrationUseCase interface {
	Handle(ctx context.Context, cmd commands.ConfirmRegistrationCommand) error
}

type SignInUseCase interface {
	Handle(ctx context.Context, cmd commands.SignInCommand) (*commands.SignInResult, error)
}

type AuthController struct {
	initiateUseCase InitiateRegistrationUseCase
	confirmUseCase  ConfirmRegistrationUseCase
	signInUseCase   SignInUseCase
}

func NewAuthController(initiate InitiateRegistrationUseCase, confirm ConfirmRegistrationUseCase, signIn SignInUseCase) *AuthController {
	return &AuthController{
		initiateUseCase: initiate,
		confirmUseCase:  confirm,
		signInUseCase:   signIn,
	}
}