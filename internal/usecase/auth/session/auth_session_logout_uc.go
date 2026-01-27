package session

import (
	"auth-service/internal/apperrs"
	"auth-service/internal/usecase/auth/credential"
	"auth-service/internal/usecase/auth/ports"
	"context"
)

type sessionLogoutUsecase struct {
	verifyCredential *credential.CredentialVerifier
	sessionRepo      ports.SessionRepository
}

func (uc *sessionLogoutUsecase) Logout(ctx context.Context, sessionID string) error {
	if err := uc.sessionRepo.Revoke(ctx, sessionID); err != nil {
		return apperrs.Internal(err)
	}

	return nil
}
