package auth

import (
	"auth-service/internal/domain/accountmodel"
	"auth-service/internal/domain/authmodel"
	"context"
	"time"
)

type TokenService interface {
	Sign(claims authmodel.TokenClaims) (string, error)
	SignPairs(jti string, signAt time.Time, account *accountmodel.Account) (*authmodel.TokenPairs, error)

	VerifyRefreshToken(ctx context.Context, token string) (*authmodel.TokenClaims, error)
}
