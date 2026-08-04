package repository

import (
	"context"

	"github.com/DimasFatchurroziq/ticket-booking-system/internal/user/domain"
	"github.com/DimasFatchurroziq/ticket-booking-system/internal/user/repository/db"
	"go.uber.org/zap"

	"github.com/google/uuid"
)

type PostgresUserRepository struct {
	queries *db.Queries
	Log     *zap.Logger
}

func NewPostgresUserRepository(conn db.DBTX, log *zap.Logger) domain.UserRepository {
	return &PostgresUserRepository{
		queries: db.New(conn),
		Log:     log,
	}
}

func toDomainUser(u db.User) *domain.User {
	return &domain.User{
		ID:           u.ID,
		PasswordHash: u.PasswordHash,
		FullName:     u.FullName,
		PhoneNumber:  u.PhoneNumber,
		CreatedAt:    u.CreatedAt.Time,
	}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	arg := db.CreateUserParams{
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		FullName:     user.FullName,
		PhoneNumber:  user.PhoneNumber,
	}

	res, err := r.queries.CreateUser(ctx, arg)
	if err != nil {
		return nil, err
	}

	return toDomainUser(res), nil
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	res, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toDomainUser(res), nil
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	res, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return toDomainUser(res), nil
}

func (r *PostgresUserRepository) List(ctx context.Context) ([]*domain.User, error) {
	rows, err := r.queries.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	var users []*domain.User
	for _, row := range rows {
		users = append(users, toDomainUser(row))
	}
	return users, nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	arg := db.UpdateUserParams{
		ID:          user.ID,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}

	res, err := r.queries.UpdateUser(ctx, arg)
	if err != nil {
		return nil, err
	}
	return toDomainUser(res), nil
}

func (r *PostgresUserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	arg := db.UpdateUserPasswordParams{
		ID:           id,
		PasswordHash: passwordHash,
	}
	return r.queries.UpdateUserPassword(ctx, arg)
}

func (r *PostgresUserRepository) UpdateEmail(ctx context.Context, id uuid.UUID, email string) (*domain.User, error) {
	arg := db.UpdateUserEmailParams{
		ID:    id,
		Email: email,
	}
	res, err := r.queries.UpdateUserEmail(ctx, arg)
	if err != nil {
		return nil, err
	}
	return toDomainUser(res), nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteUser(ctx, id)
}

func (r *PostgresUserRepository) Count(ctx context.Context) (int64, error) {
	return r.queries.CountUsers(ctx)
}

func (r *PostgresUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	exists, err := r.queries.ExistsUserByEmail(ctx, email)
	if err != nil {
		return false, err
	}
	return exists, nil
}
