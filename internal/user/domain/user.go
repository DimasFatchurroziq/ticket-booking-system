package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// User adalah Entity utama di domain user
type User struct {
	ID           uuid.UUID // Gunakan UUID untuk ID unik
	Email        string
	PasswordHash string
	FullName     string
	PhoneNumber  string
	CreatedAt    time.Time
}

// NewUser adalah aturan bisnis (factory) untuk membuat user baru dengan validasi sederhana
func NewUser(email, passwordHash, fullName, phoneNumber string) (*User, error) {
	if email == "" {
		return nil, errors.New("email tidak boleh kosong")
	}
	if passwordHash == "" {
		return nil, errors.New("password wajib di-hash")
	}

	return &User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: passwordHash,
		FullName:     fullName,
		PhoneNumber:  phoneNumber,
		CreatedAt:    time.Now(),
	}, nil
}
