package usecase

import (
	"context"

	"github.com/DimasFatchurroziq/ticket-booking-system/internal/auth/crypto"
	"github.com/DimasFatchurroziq/ticket-booking-system/internal/user/domain"
	"go.uber.org/zap"
)

type userUsecase struct {
	userRepo domain.UserRepository
	hasher   crypto.Hasher
	log      *zap.Logger
}

func NewUserUsecase(
	userRepo domain.UserRepository,
	hasher crypto.Hasher,
	log *zap.Logger,
) UserUsecase {

	return &userUsecase{
		userRepo: userRepo,
		hasher:   hasher,
		log:      log,
	}
}

func (u *userUsecase) Me(
	ctx context.Context,
	cmd MeCommand,
) (*MeResult, error) {

	u.log.Info(
		"user profile requested",
		zap.String("user_id", cmd.UserID.String()),
	)

	user, err := u.userRepo.GetByID(ctx, cmd.UserID)
	if err != nil {

		u.log.Error(
			"failed checking existing user",
			zap.String("user_id", cmd.UserID.String()),
			zap.Error(err),
		)

		return nil, err
	}

	u.log.Info(
		"user profile retrieved",
		zap.String("user_id", user.ID.String()),
		zap.String("email", user.Email),
	)

	return &MeResult{
		ID:          user.ID,
		Email:       user.Email,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}, nil
}
