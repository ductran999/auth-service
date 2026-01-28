package session

import (
	"auth-service/internal/apperrs"
	"auth-service/internal/biz/ports/auth"
	"auth-service/internal/biz/ports/repositories"
	"auth-service/internal/biz/usecase/auth/credential"
	"auth-service/internal/domain/sessionmodel"
	"context"

	"github.com/google/uuid"
)

// LoginInput represents the input required to authenticate a user using email and password.
// It also includes optional request metadata for logging, auditing, or session management.
type LoginInput struct {
	// CurrentSessionID is an optional field used to associate or trace login attempts
	CurrentSessionID string

	Email     string `json:"email"`    // User's email address
	Password  string `json:"password"` // Plain-text password from the login form
	IP        string `json:"-"`        // Client IP address (injected by handler, not from JSON)
	UserAgent string `json:"-"`        // User-Agent header string (injected by handler)
}

type sessionLoginUsecase struct {
	verifyCredential *credential.CredentialVerifier
	sessionStore     auth.SessionStore
	sessionRepo      repositories.SessionRepository
}

func (uc *sessionLoginUsecase) Login(ctx context.Context, input LoginInput) (*sessionmodel.Session, error) {
	// Verify user credentials
	account, err := uc.verifyCredential.Verify(ctx, input.Email, input.Password)
	if err != nil {
		return nil, apperrs.Unauthorized(err)
	}

	// Try to reuse existing session
	// session, err := uc.sessionStore.Refresh(ctx, input.CurrentSessionID)
	// if err != nil {
	// 	return nil, err
	// }
	// if session != nil {
	// 	return session, nil
	// }

	// Create a new session
	newSession := &sessionmodel.Session{
		ID:        uuid.New(),
		AccountID: account.ID,
		IPAddress: input.IP,
		UserAgent: input.UserAgent,
	}

	if err := uc.sessionRepo.Create(ctx, newSession); err != nil {
		return nil, apperrs.Internal(err)
	}

	return newSession, nil
}
