package usecase

import (
	"context"

	"github.com/DimasFatchurroziq/ticket-booking-system/internal/user/domain"
	"go.uber.org/zap"
)

type userUsecase struct {
	userRepo domain.UserRepository
	Log      *zap.Logger
}

func NewUserUsecase(
	userRepo domain.UserRepository,
	Log *zap.Logger,
) UserUsecase {

	return &userUsecase{
		userRepo: userRepo,
		Log:      Log,
	}
}

func (u *userUsecase) Register(
	ctx context.Context,
	cmd RegisterCommand,
) (*RegisterResult, error) {

	u.Log.Info(
		"user registration started",
		zap.String("email", cmd.Email),
	)

	user, err := domain.NewUser(
		cmd.Email,
		cmd.Password,
		cmd.FullName,
		cmd.PhoneNumber,
	)

	if err != nil {

		u.Log.Warn(
			"failed creating user entity",
			zap.String("email", cmd.Email),
			zap.Error(err),
		)

		return nil, err
	}

	exists, err := u.userRepo.ExistsByEmail(ctx, user.Email)

	if err != nil {

		u.Log.Error(
			"failed checking existing user",
			zap.String("email", user.Email),
			zap.Error(err),
		)

		return nil, err
	}

	if exists {

		u.Log.Warn(
			"user registration rejected, email already exists",
			zap.String("email", user.Email),
		)

		return nil, domain.ErrEmailExists
	}

	user, err = u.userRepo.Create(ctx, user)

	if err != nil {

		u.Log.Error(
			"failed creating user",
			zap.String("email", user.Email),
			zap.Error(err),
		)

		return nil, err
	}

	u.Log.Info(
		"user registration completed",
		zap.String("user_id", user.ID.String()),
		zap.String("email", user.Email),
	)

	return &RegisterResult{
		ID:    user.ID,
		Email: user.Email,
	}, nil
}
