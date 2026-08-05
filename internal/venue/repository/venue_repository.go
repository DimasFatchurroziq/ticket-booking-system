package repository

import (
	"context"

	"github.com/DimasFatchurroziq/ticket-booking-system/internal/venue/domain"
	"github.com/DimasFatchurroziq/ticket-booking-system/internal/venue/repository/db"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PostgresVenueRepository struct {
	queries *db.Queries
	Log     *zap.Logger
}

func NewPostgresVenueRepository(conn db.DBTX, log *zap.Logger) domain.VenueRepository {
	return &PostgresVenueRepository{
		queries: db.New(conn),
		Log:     log,
	}
}

func toDomainVenue(v db.Venue) *domain.Venue {
	return &domain.Venue{
		ID:            v.ID,
		Name:          v.Name,
		Address:       v.Address,
		City:          v.City,
		TotalCapacity: v.TotalCapacity,
		CreatedAt:     v.CreatedAt.Time,
	}
}

func (r *PostgresVenueRepository) Create(ctx context.Context, venue *domain.Venue) (*domain.Venue, error) {
	arg := db.CreateVenueParams{
		Name:          venue.Name,
		Address:       venue.Address,
		City:          venue.City,
		TotalCapacity: venue.TotalCapacity,
	}

	res, err := r.queries.CreateVenue(ctx, arg)
	if err != nil {
		return nil, err
	}

	return toDomainVenue(res), nil
}

func (r *PostgresVenueRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Venue, error) {
	res, err := r.queries.GetVenueById(ctx, id)
	if err != nil {
		return nil, err
	}
	return toDomainVenue(res), nil
}

func (r *PostgresVenueRepository) List(ctx context.Context) ([]*domain.Venue, error) {
	rows, err := r.queries.ListVenues(ctx)
	if err != nil {
		return nil, err
	}

	var venues []*domain.Venue
	for _, row := range rows {
		venues = append(venues, toDomainVenue(row))
	}
	return venues, nil
}

func (r *PostgresVenueRepository) Update(ctx context.Context, venue *domain.Venue) (*domain.Venue, error) {
	arg := db.UpdateVenueParams{
		ID:            venue.ID,
		Name:          venue.Name,
		Address:       venue.Address,
		City:          venue.City,
		TotalCapacity: venue.TotalCapacity,
	}

	res, err := r.queries.UpdateVenue(ctx, arg)
	if err != nil {
		return nil, err
	}
	return toDomainVenue(res), nil
}

func (r *PostgresVenueRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteVenue(ctx, id)
}

func (r *PostgresVenueRepository) Count(ctx context.Context) (int64, error) {
	return r.queries.CountVenues(ctx)
}
