package usecase

import "context"

type UserUsecase interface {
	Register(
		ctx context.Context,
		cmd RegisterCommand,
	) (*RegisterResult, error)
}
