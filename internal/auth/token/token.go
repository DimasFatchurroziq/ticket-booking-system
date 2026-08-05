package token

import (
	"github.com/google/uuid"
)

type TokenManager interface {
	Generate(
		userID uuid.UUID,
		email string,
	) (string, error)

	Parse(
		token string,
	) (*Claims, error)
}
