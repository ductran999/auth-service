package ports

import (
	"auth-service/internal/model"
	"context"
)

type SessionStore interface {
	Get(ctx context.Context, id string) (*model.Session, error)
	Save(ctx context.Context, s *model.Session) error
	Refresh(ctx context.Context, id string) (*model.Session, error)
}
