package session

import (
	"auth-service/internal/biz/ports/auth"
	"auth-service/internal/domain/sessionmodel"
	"context"
)

type sessionStore struct{}

func NewSessionStore() auth.SessionStore {
	return &sessionStore{}
}

func (s *sessionStore) Get(ctx context.Context, id string) (*sessionmodel.Session, error) {
	return nil, nil
}

func (s *sessionStore) Save(ctx context.Context, session *sessionmodel.Session) error {
	return nil
}

func (s *sessionStore) Refresh(ctx context.Context, id string) (*sessionmodel.Session, error) {
	return nil, nil
}
