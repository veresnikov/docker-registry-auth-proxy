package auth

import "crypto/sha256"

type PasswordHasher interface {
	Hash(password string) ([]byte, error)
}

func NewPasswordHasher() PasswordHasher {
	return &passwordHasher{}
}

type passwordHasher struct{}

func (p passwordHasher) Hash(password string) ([]byte, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(password))
	if err != nil {
		return nil, err
	}
	return hash.Sum(nil), nil
}
