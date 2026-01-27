package session

import (
	"auth-service/internal/apperrs"
	"auth-service/internal/model"
	"auth-service/internal/usecase/auth/credential"
	"auth-service/internal/usecase/auth/ports"
	"auth-service/internal/usecase/dto"
	"context"

	"github.com/google/uuid"
)

type sessionLoginUsecase struct {
	verifyCredential *credential.CredentialVerifier
	sessionStore     ports.SessionStore
	sessionRepo      ports.SessionRepository
}

func (uc *sessionLoginUsecase) Login(ctx context.Context, input dto.LoginInput) (*model.Session, error) {
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
	newSession := &model.Session{
		ID:        uuid.New(),
		AccountID: account.ID,
		Account: model.Account{
			ID:       account.ID,
			Email:    account.Email,
			Role:     account.Role,
			IsActive: account.IsActive,
		},
		IPAddress: input.IP,
		UserAgent: input.UserAgent,
	}

	if err := uc.sessionRepo.Create(ctx, newSession); err != nil {
		return nil, apperrs.Internal(err)
	}

	return newSession, nil
}
