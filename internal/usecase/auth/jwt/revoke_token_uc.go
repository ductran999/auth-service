package jwt

import (
	"auth-service/internal/usecase/auth/ports"
	"context"
)

type revokeTokenUsecase struct {
	tokenService ports.TokenService
	tokenStore   ports.TokenStore
}

func (uc *revokeTokenUsecase) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return ErrInvalidRefreshToken
	}

	claims, err := uc.tokenService.VerifyRefreshToken(ctx, refreshToken)
	if err != nil {
		return ErrInvalidRefreshToken
	}

	if err := uc.tokenStore.Revoke(ctx, claims.Subject, claims.ID); err != nil {
		return err
	}

	return nil
}
