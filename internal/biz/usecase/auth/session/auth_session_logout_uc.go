package session

import (
	"auth-service/internal/apperrs"

	"auth-service/internal/biz/ports/repositories"
	"auth-service/internal/biz/usecase/auth/credential"
	"context"
)

type sessionLogoutUsecase struct {
	verifyCredential *credential.CredentialVerifier
	sessionRepo      repositories.SessionRepository
}

func (uc *sessionLogoutUsecase) Logout(ctx context.Context, sessionID string) error {
	if err := uc.sessionRepo.Revoke(ctx, sessionID); err != nil {
		return apperrs.Internal(err)
	}

	return nil
}
