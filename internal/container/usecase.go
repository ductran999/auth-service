package container

import (
	accountUsecase "auth-service/internal/usecase/account"
	"auth-service/internal/usecase/auth/credential"
	authJwtUsecase "auth-service/internal/usecase/auth/jwt"
	authSessionUsecase "auth-service/internal/usecase/auth/session"

	"auth-service/internal/infra/storage/redis"
	"auth-service/internal/infra/token/signer"
	"auth-service/internal/usecase/port"
)

type useCases struct {
	jwtAuth     port.AuthJWTUsecase
	sessionAuth port.AuthSessionUsecase

	account port.AccountUsecase
	session port.SessionUsecase

	backgroundSession port.SessionMaintenanceUsecase
}

func (c *Container) initUseCases() {
	accountUC := accountUsecase.NewAccountUseCase(
		c.Hasher,
		c.repositories.account,
	)

	// sessionAuthUC := sessionUsecase.NewAuthSessionUsecase(c.Cache, c.repositories.session, c.Hasher, c.repositories.account)

	tokenService := signer.NewJWTSigner(c.Signer)
	tokenStore := redis.NewTokenStore(c.Cache)
	credVerifier := credential.NewCredentialVerifier(c.Hasher, c.repositories.account)

	jwtAuthUC := authJwtUsecase.NewAuthJWTUsecase(tokenService, tokenStore, credVerifier)

	sessionStore := redis.NewSessionStore()
	sessionAuthUC := authSessionUsecase.NewAuthSessionUsecase(credVerifier, sessionStore, c.repositories.session)

	c.useCases = &useCases{
		account:     accountUC,
		jwtAuth:     jwtAuthUC,
		sessionAuth: sessionAuthUC,
		// session:           sessionUC,
		// backgroundSession: sessionUC,
	}
}
