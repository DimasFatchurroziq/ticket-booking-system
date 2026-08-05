package usecase

import (
	"context"

	"github.com/google/uuid"
)

type AuthUsecase interface {
	Register(
		ctx context.Context,
		cmd RegisterCommand,
	) (*RegisterResult, error)
	Login(
		ctx context.Context,
		cmd LoginCommand,
	) (*LoginResult, error)
}

type RegisterCommand struct {
	Email       string
	Password    string
	FullName    string
	PhoneNumber string
}

type RegisterResult struct {
	ID    uuid.UUID
	Email string
}

type LoginCommand struct {
	Email    string
	Password string
}

type LoginResult struct {
	AccessToken string
}
