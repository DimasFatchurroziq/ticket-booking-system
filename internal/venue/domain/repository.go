package domain

import (
	"context"

	"github.com/google/uuid"
)

type VenueRepository interface {
	Create(ctx context.Context, venue *Venue) (*Venue, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Venue, error)
	List(ctx context.Context) ([]*Venue, error)
	Update(ctx context.Context, venue *Venue) (*Venue, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Count(ctx context.Context) (int64, error)
}
