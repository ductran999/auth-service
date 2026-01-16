package port

import (
	"context"

	"auth-service/internal/usecase/dto"
)

type AuthJWTUsecase interface {
	Login(ctx context.Context, input dto.LoginJWTInput) (*dto.TokenPairs, error)

	RefreshToken(ctx context.Context, refreshToken string) (*dto.TokenPairs, error)

	RevokeRefreshToken(ctx context.Context, refreshToken string) error
}
