package jwt

import (
	"auth-service/internal/usecase/auth/credential"
	"auth-service/internal/usecase/auth/ports"
	"auth-service/internal/usecase/port"
)

type authJWT struct {
	*loginWithJWTUsecase
	*refreshTokenUsecase
	*revokeTokenUsecase
}

func NewAuthJWTUsecase(
	tokenService ports.TokenService,
	tokenStore ports.TokenStore,
	credentialVerifier *credential.CredentialVerifier,
) port.AuthJWTUsecase {
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
