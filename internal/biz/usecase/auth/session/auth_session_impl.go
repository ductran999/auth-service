package session

import (
	"context"
	"time"

	"auth-service/internal/biz/ports/auth"
	"auth-service/internal/biz/ports/repositories"
	"auth-service/internal/biz/usecase/auth/credential"
	"auth-service/internal/domain/authmodel"
	"auth-service/internal/domain/sessionmodel"
)

const (
	SessionLifetime = 60 * time.Minute
)

// AuthSessionUsecase defines the authentication-related business logic based on user sessions.
type AuthSessionUsecase interface {
	// Login verifies the provided credentials and returns the authenticated account.
	// Returns an error if authentication fails.
	Login(ctx context.Context, input LoginInput) (*sessionmodel.Session, error)

	// Logout terminates the session associated with the given session ID.
	// It removes the session from cache (best-effort) and marks it as expired in the database.
	// Returns an error only if the database update fails.
	Logout(ctx context.Context, sessionID string) error

	ValidateSession(ctx context.Context, sessionID string) (*authmodel.AuthObj, error)
}

type authSessionUsecase struct {
	*sessionLoginUsecase
	*sessionLogoutUsecase
	*validateSessionUsecase
}

func NewAuthSessionUsecase(
	verifyCredential *credential.CredentialVerifier,
	sessionStore auth.SessionStore,
	sessionRepo repositories.SessionRepository,
) AuthSessionUsecase {
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
		validateSessionUsecase: &validateSessionUsecase{
			sessionStore: sessionStore,
		},
	}
}
