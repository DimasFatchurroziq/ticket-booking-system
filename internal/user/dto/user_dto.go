package dto

import "github.com/google/uuid"

type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	FullName    string `json:"full_name" validate:"required,min=3"`
	PhoneNumber string `json:"phone_number" validate:"required,numeric"`
}

type RegisterResponse struct {
	ID      uuid.UUID `json:"id"`
	Email   string    `json:"email"`
	Message string    `json:"message"`
}
