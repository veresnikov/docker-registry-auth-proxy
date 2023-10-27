package storage

import (
	"bytes"
	"context"

	"github.com/veresnikov/docker-registry-auth-proxy/pkg/application/auth"
	"github.com/veresnikov/docker-registry-auth-proxy/pkg/infrastructure/config"
)

func NewUserStorage(configPath string, passwordHasher auth.PasswordHasher) (auth.UserStorage, error) {
	data, err := config.Load(configPath, passwordHasher)
	if err != nil {
		return nil, err
	}
	return &userStorage{
		data: data,
	}, nil
}

type userStorage struct {
	data map[config.UserName]config.PasswordHash
}

func (u userStorage) Authorize(_ context.Context, username string, password []byte) (bool, error) {
	p := u.data[config.UserName(username)]
	return bytes.Equal(password, p), nil
}
