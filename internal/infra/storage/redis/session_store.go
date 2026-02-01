package redis

import (
	"auth-service/internal/model"
	"auth-service/internal/usecase/auth/ports"
	"context"
)

type sessionStore struct{}

func NewSessionStore() ports.SessionStore {
	return &sessionStore{}
}

func (s *sessionStore) Get(ctx context.Context, id string) (*model.Session, error) {
	return nil, nil
}

func (s *sessionStore) Save(ctx context.Context, session *model.Session) error {
	return nil
}

func (s *sessionStore) Refresh(ctx context.Context, id string) (*model.Session, error) {
	return nil, nil
}
