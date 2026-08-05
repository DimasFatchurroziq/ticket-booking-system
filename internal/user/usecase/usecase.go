package usecase

import (
	"context"

	"github.com/google/uuid"
)

type UserUsecase interface {
	Me(
		ctx context.Context,
		cmd MeCommand,
	) (*MeResult, error)
}

type MeCommand struct {
	UserID uuid.UUID
}

type MeResult struct {
	ID          uuid.UUID
	Email       string
	FullName    string
	PhoneNumber string
}
