package jwt

import (
	"auth-service/internal/apperrs"
	"auth-service/internal/biz/ports/auth"
	"auth-service/internal/domain/authmodel"
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type refreshTokenUsecase struct {
	tokenService auth.TokenService
	tokenStore   auth.TokenStore
}

func (uc *refreshTokenUsecase) RefreshToken(ctx context.Context, refreshToken string) (*authmodel.TokenPairs, error) {
	if refreshToken == "" {
		return nil, apperrs.Unauthorized(ErrInvalidRefreshToken)
	}

	claims, err := uc.tokenService.VerifyRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, apperrs.Unauthorized(ErrInvalidRefreshToken)
	}

	ok, err := uc.tokenStore.Exists(ctx, claims.ID)
	if err != nil {
		return nil, apperrs.Internal(err)
	}
	if !ok {
		return nil, apperrs.Unauthorized(ErrInvalidRefreshToken)
	}

	tokens, err := uc.resignTokenPairs(ctx, *claims)
	if err != nil {
		return nil, apperrs.Internal(err)
	}

	return tokens, nil
}

func (uc *refreshTokenUsecase) resignTokenPairs(ctx context.Context, old authmodel.TokenClaims) (*authmodel.TokenPairs, error) {
	now := time.Now()

	// Build access token claims (no JTI)
	accessClaims := authmodel.TokenClaims{
		Email: old.Email,
		Role:  old.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   old.Subject,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(AccessTokenLifetime)),
		},
	}
	accessToken, err := uc.tokenService.Sign(accessClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	//  Build refresh token claims (new JTI)
	jti := uuid.NewString()
	refreshClaims := authmodel.TokenClaims{
		Email: old.Email,
		Role:  old.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			Subject:   old.Subject,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(RefreshTokenLifetime)),
		},
	}
	refreshToken, err := uc.tokenService.Sign(refreshClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	deviceSession := authmodel.DeviceSession{
		JTI:       jti,
		AccountID: old.Subject,
		SignAt:    now,
		ExpiresAt: now.Add(RefreshTokenLifetime),
	}
	if err := uc.tokenStore.Save(ctx, deviceSession); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	// Revoke old refresh token
	if err := uc.tokenStore.Revoke(ctx, old.Subject, old.ID); err != nil {
		return nil, fmt.Errorf("failed to invalidate old refresh token: %w", err)
	}

	return &authmodel.TokenPairs{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
