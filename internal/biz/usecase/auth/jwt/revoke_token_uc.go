package jwt

import (
	"auth-service/internal/apperrs"
	"auth-service/internal/biz/ports/auth"
	"context"
)

type revokeTokenUsecase struct {
	tokenService auth.TokenService
	tokenStore   auth.TokenStore
}

func (uc *revokeTokenUsecase) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return apperrs.Unauthorized(ErrInvalidRefreshToken)
	}

	claims, err := uc.tokenService.VerifyRefreshToken(ctx, refreshToken)
	if err != nil {
		return apperrs.Unauthorized(ErrInvalidRefreshToken)
	}

	_ = uc.tokenStore.Revoke(ctx, claims.Subject, claims.ID)

	return nil
}
