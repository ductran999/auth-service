package session

import (
	"auth-service/internal/model"
	sessionUC "auth-service/internal/usecase/session"
	"context"
	"errors"

	"gorm.io/gorm"
)

func (r *sessionRepoImpl) FindByID(ctx context.Context, sessionID string) (*model.Session, error) {
	var session model.Session

	err := r.db.WithContext(ctx).
		Preload("Account", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "email", "role", "is_active")
		}).
		Where("id = ?", sessionID).
		First(&session).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sessionUC.ErrSessionNotFound
		}
		return nil, err
	}

	return &session, nil
}

func (r *sessionRepoImpl) FindAllActiveSession(ctx context.Context) ([]model.Session, error) {
	var activeSessions []model.Session

	err := r.db.WithContext(ctx).
		Select("id").
		Where("expires_at IS NULL").
		Find(&activeSessions).Error

	if err != nil {
		return nil, err
	}

	return activeSessions, nil
}
