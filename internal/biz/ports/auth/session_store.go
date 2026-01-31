package auth

import (
	"auth-service/internal/domain/sessionmodel"
	"context"
)

type SessionStore interface {
	Get(ctx context.Context, id string) (*sessionmodel.Session, error)
	Save(ctx context.Context, s *sessionmodel.Session) error
	Refresh(ctx context.Context, id string) (*sessionmodel.Session, error)
}
