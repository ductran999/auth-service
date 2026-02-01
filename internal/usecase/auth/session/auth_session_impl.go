package session

import (
	"context"
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
	}
}

func (uc *authSessionUsecase) Logout(ctx context.Context, sessionID string) error {
	// Remove session in cache
	// _ = uc.cache.Del(ctx, cache.KeyFromSessionID(sessionID))

	// expire the session in db
	// err := uc.sessionRepo.UpdateExpiresAt(ctx, sessionID, time.Now())
	// if err != nil {
	// return fmt.Errorf("logout: failed to update expires_at in DB for session=%s, error=%w", sessionID, err)
	// }

	return nil
}

// func (uc *authSessionUsecase) createSession(
// 	ctx context.Context,
// 	account *model.Account,
// 	input dto.LoginInput,
// ) (*model.Session, error) {
// 	session := &model.Session{
// 		AccountID: account.ID,
// 		Account: model.Account{
// 			ID:       account.ID,
// 			Email:    account.Email,
// 			Role:     account.Role,
// 			IsActive: account.IsActive,
// 		},
// 		IPAddress: input.IP,
// 		UserAgent: input.UserAgent,
// 	}

// 	if err := uc.sessionRepo.Create(ctx, session); err != nil {
// 		return nil, err
// 	}

// 	_ = uc.cache.Set(ctx, cache.KeyFromSessionID(session.ID.String()), session, SessionLifetime)

// 	return session, nil
// }

// func (uc *authSessionUsecase) getSessionFromCache(
// 	ctx context.Context,
// 	sessionKey string,
// ) *model.Session {
// 	var session model.Session
// 	if err := uc.cache.GetInto(ctx, sessionKey, &session); err != nil {
// 		// If get session from cache got error just pass it.
// 		// Already has fallback form DB
// 		return nil
// 	}

// 	return &session
// }
