package auth

import (
	"auth-service/internal/domain/authmodel"
	"context"
)

type TokenStore interface {
	Save(ctx context.Context, device authmodel.DeviceSession) error
	Revoke(ctx context.Context, userID, tokenID string) error

	Exists(ctx context.Context, userID string, tokenID string) (bool, error)
}
