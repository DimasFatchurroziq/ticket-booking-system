package usecase

import "github.com/google/uuid"

type RegisterResult struct {
	ID    uuid.UUID
	Email string
}
