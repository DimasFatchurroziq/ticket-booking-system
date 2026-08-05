package domain

import (
	"context"

	"github.com/google/uuid"
)

type EventRepository interface {
	Create(ctx context.Context, event *Event) (*Event, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Event, error)
	List(ctx context.Context, filter *EventFilter) ([]*Event, error)
	Update(ctx context.Context, event *Event) (*Event, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Count(ctx context.Context) (int64, error)
}
