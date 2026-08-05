package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTManager struct {
	secretKey []byte
	issuer    string
	expiredAt time.Duration
}

func NewJWTManager(
	secret string,
	issuer string,
	expiredAt time.Duration,
) *JWTManager {

	return &JWTManager{
		secretKey: []byte(secret),
		issuer:    issuer,
		expiredAt: expiredAt,
	}
}

func (j *JWTManager) Generate(
	userID uuid.UUID,
	email string,
) (string, error) {

	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expiredAt)),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(j.secretKey)
}

func (j *JWTManager) Parse(
	tokenString string,
) (*Claims, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return j.secretKey, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
