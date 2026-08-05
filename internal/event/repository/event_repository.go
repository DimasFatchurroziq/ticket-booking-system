package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"

	"github.com/DimasFatchurroziq/ticket-booking-system/internal/event/domain"
	"github.com/DimasFatchurroziq/ticket-booking-system/internal/event/repository/db"
)

type PostgresEventRepository struct {
	queries *db.Queries
	Log     *zap.Logger
}

func NewPostgresEventRepository(conn db.DBTX, log *zap.Logger) domain.EventRepository {
	return &PostgresEventRepository{
		queries: db.New(conn),
		Log:     log,
	}
}

func toPgTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  t,
		Valid: true,
	}
}

func toDomainEvent(e db.Event) *domain.Event {
	return &domain.Event{
		ID:          e.ID,
		VenueID:     e.VenueID,
		Name:        e.Name,
		Description: e.Description,

		EventStart: e.EventStart.Time,
		EventEnd:   e.EventEnd.Time,

		SaleStart: e.SaleStart.Time,
		SaleEnd:   e.SaleEnd.Time,

		Status: domain.EventStatus(e.Status),

		CreatedAt: e.CreatedAt.Time,
		UpdatedAt: e.UpdatedAt.Time,
	}
}

func (r *PostgresEventRepository) Create(ctx context.Context, event *domain.Event) (*domain.Event, error) {

	arg := db.CreateEventParams{
		VenueID:     event.VenueID,
		Name:        event.Name,
		Description: event.Description,

		EventStart: toPgTimestamptz(event.EventStart),
		EventEnd:   toPgTimestamptz(event.EventEnd),

		SaleStart: toPgTimestamptz(event.SaleStart),
		SaleEnd:   toPgTimestamptz(event.SaleEnd),

		Status: string(event.Status),
	}

	res, err := r.queries.CreateEvent(ctx, arg)
	if err != nil {
		return nil, err
	}

	return toDomainEvent(res), nil
}

func (r *PostgresEventRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {

	res, err := r.queries.GetEvent(ctx, id)
	if err != nil {
		return nil, err
	}

	return toDomainEvent(res), nil
}

func (r *PostgresEventRepository) List(ctx context.Context, filter *domain.EventFilter) ([]*domain.Event, error) {

	if filter == nil {
		filter = &domain.EventFilter{}
	}

	params := db.ListEventsParams{
		FilterStatus:        pgtype.Text{},
		FilterStartSaleDate: pgtype.Timestamptz{},
		FilterEndSaleDate:   pgtype.Timestamptz{},
	}

	if filter.Status != nil {
		params.FilterStatus = pgtype.Text{
			String: string(*filter.Status),
			Valid:  true,
		}
	}

	if filter.StartSaleDate != nil {
		params.FilterStartSaleDate = pgtype.Timestamptz{
			Time:  *filter.StartSaleDate,
			Valid: true,
		}
	}

	if filter.EndSaleDate != nil {
		params.FilterEndSaleDate = pgtype.Timestamptz{
			Time:  *filter.EndSaleDate,
			Valid: true,
		}
	}

	rows, err := r.queries.ListEvents(ctx, params)
	if err != nil {
		return nil, err
	}

	var events []*domain.Event

	for _, row := range rows {
		events = append(events, toDomainEvent(row))
	}

	return events, nil
}

func (r *PostgresEventRepository) Update(ctx context.Context, event *domain.Event) (*domain.Event, error) {

	arg := db.UpdateEventParams{
		ID:          event.ID,
		VenueID:     event.VenueID,
		Name:        event.Name,
		Description: event.Description,

		EventStart: toPgTimestamptz(event.EventStart),
		EventEnd:   toPgTimestamptz(event.EventEnd),

		SaleStart: toPgTimestamptz(event.SaleStart),
		SaleEnd:   toPgTimestamptz(event.SaleEnd),

		Status: string(event.Status),
	}

	res, err := r.queries.UpdateEvent(ctx, arg)
	if err != nil {
		return nil, err
	}

	return toDomainEvent(res), nil
}

func (r *PostgresEventRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteEvent(ctx, id)
}

func (r *PostgresEventRepository) Count(ctx context.Context) (int64, error) {
	return r.queries.CountEvents(ctx)
}
