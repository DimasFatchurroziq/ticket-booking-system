package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Venue struct {
	ID            uuid.UUID // Gunakan UUID untuk ID unik
	Name          string
	Address       string
	City          string
	TotalCapacity int32
	CreatedAt     time.Time
}

func NewVenue(name, address, city string, totalCapacity int32) (*Venue, error) {
	if name == "" {
		return nil, errors.New("name tidak boleh kosong")
	}
	if address == "" {
		return nil, errors.New("address tidak boleh kosong")
	}
	if city == "" {
		return nil, errors.New("city tidak boleh kosong")
	}
	if totalCapacity <= 0 {
		return nil, errors.New("total capacity tidak boleh nol atau negatif")
	}

	return &Venue{
		ID:            uuid.New(),
		Name:          name,
		Address:       address,
		City:          city,
		TotalCapacity: totalCapacity,
		CreatedAt:     time.Now(),
	}, nil
}
