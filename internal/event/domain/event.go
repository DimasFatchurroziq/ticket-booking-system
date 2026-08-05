package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type EventStatus string

const (
	EventStatusDraft     EventStatus = "draft"
	EventStatusPublished EventStatus = "published"
	EventStatusOngoing   EventStatus = "ongoing"
	EventStatusCompleted EventStatus = "completed"
	EventStatusCancelled EventStatus = "cancelled"
)

type Event struct {
	ID          uuid.UUID
	VenueID     uuid.UUID
	Name        string
	Description string

	EventStart time.Time
	EventEnd   time.Time

	SaleStart time.Time
	SaleEnd   time.Time

	Status EventStatus

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewEvent(
	venueID uuid.UUID,
	name string,
	description string,
	eventStart time.Time,
	eventEnd time.Time,
	saleStart time.Time,
	saleEnd time.Time,
) (*Event, error) {

	if venueID == uuid.Nil {
		return nil, errors.New("venue id tidak valid")
	}

	if name == "" {
		return nil, errors.New("nama event tidak boleh kosong")
	}

	if eventEnd.Before(eventStart) {
		return nil, errors.New("event end harus setelah event start")
	}

	if saleEnd.Before(saleStart) {
		return nil, errors.New("sale end harus setelah sale start")
	}

	if saleEnd.After(eventStart) {
		return nil, errors.New("penjualan tiket harus selesai sebelum event dimulai")
	}

	now := time.Now()

	return &Event{
		ID:          uuid.New(),
		VenueID:     venueID,
		Name:        name,
		Description: description,
		EventStart:  eventStart,
		EventEnd:    eventEnd,
		SaleStart:   saleStart,
		SaleEnd:     saleEnd,
		Status:      EventStatusDraft,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

type EventFilter struct {
	Status *EventStatus

	StartSaleDate *time.Time
	EndSaleDate   *time.Time
}
