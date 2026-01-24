package ports

import (
	"auth-service/internal/model"
	"auth-service/internal/usecase/dto"
	"context"
	"time"
)

type TokenService interface {
	Sign(claims model.TokenClaims) (string, error)
	SignPairs(jti string, signAt time.Time, account *model.Account) (*dto.TokenPairs, error)

	VerifyAccessToken(ctx context.Context, token string) (model.TokenClaims, error)
	VerifyRefreshToken(ctx context.Context, token string) (model.TokenClaims, error)
}
