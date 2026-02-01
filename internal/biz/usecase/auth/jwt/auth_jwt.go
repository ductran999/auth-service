package jwt

import (
	"auth-service/internal/biz/ports/auth"
	"auth-service/internal/biz/usecase/auth/credential"
	"auth-service/internal/domain/authmodel"
	"context"
)

type AuthJWTUsecase interface {
	Login(ctx context.Context, input LoginJWTInput) (*authmodel.TokenPairs, error)
	RefreshToken(ctx context.Context, refreshToken string) (*authmodel.TokenPairs, error)
	RevokeRefreshToken(ctx context.Context, refreshToken string) error
}

type authJWT struct {
	*loginWithJWTUsecase
	*refreshTokenUsecase
	*revokeTokenUsecase
}

func NewAuthJWTUsecase(
	tokenService auth.TokenService,
	tokenStore auth.TokenStore,
	credentialVerifier *credential.CredentialVerifier,
) AuthJWTUsecase {
	return &authJWT{
		loginWithJWTUsecase: &loginWithJWTUsecase{
			tokenService: tokenService,
			tokenStore:   tokenStore,
			credVerifier: credentialVerifier,
		},
		refreshTokenUsecase: &refreshTokenUsecase{
			tokenService: tokenService,
			tokenStore:   tokenStore,
		},
		revokeTokenUsecase: &revokeTokenUsecase{
			tokenService: tokenService,
			tokenStore:   tokenStore,
		},
	}
}
