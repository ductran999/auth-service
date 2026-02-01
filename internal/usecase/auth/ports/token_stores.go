package ports

import (
	"auth-service/internal/model"
	"context"
)

type TokenStore interface {
	Save(ctx context.Context, device model.DeviceSession) error
	Revoke(ctx context.Context, userID, tokenID string) error

	Exists(ctx context.Context, tokenID string) (bool, error)
}
