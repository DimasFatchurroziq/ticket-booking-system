package usecase

import (
	"context"

	"github.com/DimasFatchurroziq/ticket-booking-system/internal/auth/crypto"
	"github.com/DimasFatchurroziq/ticket-booking-system/internal/auth/domain"
	"github.com/DimasFatchurroziq/ticket-booking-system/internal/auth/token"

	"go.uber.org/zap"
)

type authUsecase struct {
	userRepo     domain.AuthRepository
	hasher       crypto.Hasher
	tokenManager token.TokenManager
	log          *zap.Logger
}

func NewAuthUsecase(
	userRepo domain.AuthRepository,
	hasher crypto.Hasher,
	tokenManager token.TokenManager,
	log *zap.Logger,
) AuthUsecase {

	return &authUsecase{
		userRepo:     userRepo,
		hasher:       hasher,
		tokenManager: tokenManager,
		log:          log,
	}
}

func (u *authUsecase) Register(ctx context.Context, cmd RegisterCommand) (*RegisterResult, error) {
	u.log.Info(
		"user registration started",
		zap.String("email", cmd.Email),
	)

	exists, err := u.userRepo.ExistsByEmail(ctx, cmd.Email)
	if err != nil {

		u.log.Error(
			"failed checking existing user",
			zap.String("email", cmd.Email),
			zap.Error(err),
		)

		return nil, err
	}

	if exists {

		u.log.Warn(
			"user registration rejected, email already exists",
			zap.String("email", cmd.Email),
		)

		return nil, domain.ErrEmailExists
	}

	passwordHash, err := u.hasher.Hash(cmd.Password)
	if err != nil {

		u.log.Error(
			"failed hashing password",
			zap.String("email", cmd.Email),
			zap.Error(err),
		)

		return nil, err
	}

	user, err := domain.NewUser(
		cmd.Email,
		passwordHash,
		cmd.FullName,
		cmd.PhoneNumber,
	)
	if err != nil {

		u.log.Warn(
			"failed creating user entity",
			zap.String("email", cmd.Email),
			zap.Error(err),
		)

		return nil, err
	}

	user, err = u.userRepo.Create(ctx, user)
	if err != nil {

		u.log.Error(
			"failed creating user",
			zap.String("email", cmd.Email),
			zap.Error(err),
		)

		return nil, err
	}

	u.log.Info(
		"user registration completed",
		zap.String("user_id", user.ID.String()),
		zap.String("email", user.Email),
	)

	return &RegisterResult{
		ID:    user.ID,
		Email: user.Email,
	}, nil
}

func (u *authUsecase) Login(ctx context.Context, cmd LoginCommand) (*LoginResult, error) {
	u.log.Info(
		"user login started",
		zap.String("email", cmd.Email),
	)

	user, err := u.userRepo.GetByEmail(ctx, cmd.Email)
	if err != nil {

		u.log.Error(
			"failed checking existing user",
			zap.String("email", cmd.Email),
			zap.Error(err),
		)

		return nil, err
	}

	err = u.hasher.Compare(user.PasswordHash, cmd.Password)

	if err != nil {

		u.log.Error(
			"failed comparing password",
			zap.String("email", cmd.Email),
			zap.Error(err),
		)

		return nil, err
	}

	signedToken, err := u.tokenManager.Generate(user.ID, user.Email)
	if err != nil {
		u.log.Error(
			"failed generating token",
			zap.String("email", cmd.Email),
			zap.Error(err),
		)

		return nil, err
	}

	u.log.Info(
		"user login completed",
		zap.String("user_id", user.ID.String()),
		zap.String("email", user.Email),
	)

	return &LoginResult{
		AccessToken: signedToken,
	}, nil
}
