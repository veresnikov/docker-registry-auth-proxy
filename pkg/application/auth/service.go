package auth

import (
	"context"
	stderrors "errors"
)

var (
	ErrUnauthorized = stderrors.New("unauthorized")
)

type UserStorage interface {
	Authorize(ctx context.Context, username string, password []byte) (bool, error)
}

type Service interface {
	Authorize(ctx context.Context, username string, password []byte) error
}

func NewService(storage UserStorage) Service {
	return &service{storage: storage}
}

type service struct {
	storage UserStorage
}

func (s service) Authorize(ctx context.Context, username string, password []byte) error {
	isAuthorize, err := s.storage.Authorize(ctx, username, password)
	if err != nil {
		return err
	}
	if isAuthorize {
		return nil
	}
	return ErrUnauthorized
}
