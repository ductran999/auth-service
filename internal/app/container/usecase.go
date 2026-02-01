package container

import (
	accountUsecase "auth-service/internal/biz/usecase/account"
	"auth-service/internal/biz/usecase/auth/credential"
	authJwtUsecase "auth-service/internal/biz/usecase/auth/jwt"
	authSessionUsecase "auth-service/internal/biz/usecase/auth/session"
	"auth-service/internal/infra/auth"
	"auth-service/internal/infra/session"
)

type useCases struct {
	jwtAuth     authJwtUsecase.AuthJWTUsecase
	sessionAuth authSessionUsecase.AuthSessionUsecase

	account accountUsecase.AccountUsecase
}

func (c *Container) initUseCases() {
	accountUC := accountUsecase.NewAccountUseCase(
		c.Hasher,
		c.repos.account,
	)

	tokenService := auth.NewJWTSigner(c.Signer)
	tokenStore := auth.NewTokenStore(c.Cache)
	credVerifier := credential.NewCredentialVerifier(c.Hasher, c.repos.account)

	jwtAuthUC := authJwtUsecase.NewAuthJWTUsecase(tokenService, tokenStore, credVerifier)

	sessionStore := session.NewSessionStore(c.Cache)
	sessionAuthUC := authSessionUsecase.NewAuthSessionUsecase(credVerifier, sessionStore, c.repos.session)

	c.useCases = &useCases{
		account:     accountUC,
		jwtAuth:     jwtAuthUC,
		sessionAuth: sessionAuthUC,
	}
}

func (c *Container) GetAuthSessionUC() authSessionUsecase.AuthSessionUsecase {
	return c.useCases.sessionAuth
}
