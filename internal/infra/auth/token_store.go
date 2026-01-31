package auth

import (
	"auth-service/internal/biz/ports/auth"
	"auth-service/internal/domain/authmodel"
	"context"

	"github.com/DucTran999/cachekit"
)

const (
	RedisRefreshKeyPrefix = "auth:refresh:"
)

type tokenStore struct {
	redisClient cachekit.RemoteCache
}

func NewTokenStore(cl cachekit.RemoteCache) auth.TokenStore {
	return &tokenStore{
		redisClient: cl,
	}
}

func (ts *tokenStore) Save(ctx context.Context, device authmodel.DeviceSession) error {
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
