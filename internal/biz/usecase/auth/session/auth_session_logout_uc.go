package session

import (
	"auth-service/internal/apperrs"
<<<<<<< HEAD:internal/biz/usecase/auth/session/auth_session_logout_uc.go
	"auth-service/internal/biz/ports/repositories"
	"auth-service/internal/biz/usecase/auth/credential"
=======
	"auth-service/internal/usecase/auth/credential"
	"auth-service/internal/usecase/auth/ports"
>>>>>>> b591f14c3379b2ad171c024ba426d595c7da4137:internal/usecase/auth/session/auth_session_logout_uc.go
	"context"
)

type sessionLogoutUsecase struct {
	verifyCredential *credential.CredentialVerifier
<<<<<<< HEAD:internal/biz/usecase/auth/session/auth_session_logout_uc.go
	sessionRepo      repositories.SessionRepository
=======
	sessionRepo      ports.SessionRepository
>>>>>>> b591f14c3379b2ad171c024ba426d595c7da4137:internal/usecase/auth/session/auth_session_logout_uc.go
}

func (uc *sessionLogoutUsecase) Logout(ctx context.Context, sessionID string) error {
	if err := uc.sessionRepo.Revoke(ctx, sessionID); err != nil {
		return apperrs.Internal(err)
	}

	return nil
}
