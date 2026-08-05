package crypto

import "golang.org/x/crypto/bcrypt"

type BcryptHasher struct {
	cost int
}

func NewBcryptHasher(cost int) *BcryptHasher {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}

	return &BcryptHasher{
		cost: cost,
	}
}

func (b *BcryptHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		b.cost,
	)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (b *BcryptHasher) Compare(hash, password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
}
