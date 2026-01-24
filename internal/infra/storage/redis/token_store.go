package redis

import (
	"auth-service/internal/model"
	"auth-service/internal/usecase/auth/ports"
	"context"
	"time"

	"github.com/DucTran999/cachekit"
)

const (
	SessionKeyPrefix      = "auth:session:"
	RedisRefreshKeyPrefix = "auth:refresh:"
)

const (
	AccessTokenLifetime  = 15 * time.Minute
	RefreshTokenLifetime = 7 * 24 * time.Hour
)

type tokenStore struct {
	redisClient cachekit.RemoteCache
}

func NewTokenStore(cl cachekit.RemoteCache) ports.TokenStore {
	return &tokenStore{
		redisClient: cl,
	}
}

func (ts *tokenStore) Save(ctx context.Context, device model.DeviceSession) error {
	key := ts.keyRefreshToken(device.AccountID, device.JTI)
	if err := ts.redisClient.Set(ctx, key, device, RefreshTokenLifetime); err != nil {
		return err
	}

	return nil
}

func (ts *tokenStore) Revoke(ctx context.Context, userID, tokenID string) error {
	return nil
}

func (ts *tokenStore) Exists(ctx context.Context, tokenID string) (bool, error) {
	return false, nil
}

// KeyRefreshToken returns the Redis key for a specific refresh token.
//
//   - format: "auth:refresh:<user_id>:<jti>"
//   - example: "auth:refresh:42:b7f3-xyz"
func (ts *tokenStore) keyRefreshToken(userID, jti string) string {
	return RedisRefreshKeyPrefix + userID + ":" + jti
}
