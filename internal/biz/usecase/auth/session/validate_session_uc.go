package session

import (
	"auth-service/internal/apperrs"
	"auth-service/internal/biz/ports/auth"
	"auth-service/internal/domain/authmodel"
	"context"
	"errors"
)

var (
	ErrInvalidSession = errors.New("unauthorized session")
)

type validateSessionUsecase struct {
	sessionStore auth.SessionStore
}

func (uc *validateSessionUsecase) ValidateSession(ctx context.Context, sessionID string) (*authmodel.AuthObj, error) {
	session, err := uc.sessionStore.Get(ctx, sessionID)
	if err != nil {
		return nil, apperrs.Unauthorized(ErrInvalidSession)
	}

	authObj := authmodel.AuthObj{
		UserID: session.AccountID.String(),
	}

	return &authObj, nil
}
