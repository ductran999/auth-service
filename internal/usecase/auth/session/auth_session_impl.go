package session

import (
	"time"

	"auth-service/internal/usecase/auth/credential"
	"auth-service/internal/usecase/auth/ports"
	"auth-service/internal/usecase/port"
)

const (
	SessionLifetime = 60 * time.Minute
)

type authSessionUsecase struct {
	*sessionLoginUsecase
	*sessionLogoutUsecase
}

func NewAuthSessionUsecase(
	verifyCredential *credential.CredentialVerifier,
	sessionStore ports.SessionStore,
	sessionRepo ports.SessionRepository,
) port.AuthSessionUsecase {
	return &authSessionUsecase{
		sessionLoginUsecase: &sessionLoginUsecase{
			verifyCredential: verifyCredential,
			sessionStore:     sessionStore,
			sessionRepo:      sessionRepo,
		},
		sessionLogoutUsecase: &sessionLogoutUsecase{
			verifyCredential: verifyCredential,
			sessionRepo:      sessionRepo,
		},
	}
}
