package session

import (
	"auth-service/internal/biz/ports/repositories"
	"auth-service/internal/domain/sessionmodel"
	"context"
	"time"

	"gorm.io/gorm"
)

// sessionRepo implements the SessionRepository interface.
type sessionRepoImpl struct {
	db *gorm.DB
}

// NewSessionRepository returns a concrete implementation of SessionRepository.
func NewSessionRepository(db *gorm.DB) repositories.SessionRepository {
	return &sessionRepoImpl{db: db}
}

// Create inserts a new session record into the database.
func (r *sessionRepoImpl) Create(ctx context.Context, session *sessionmodel.Session) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *sessionRepoImpl) Revoke(ctx context.Context, sessionID string) error {
	return r.db.WithContext(ctx).
		Model(&sessionmodel.Session{}).
		Where("id = ?", sessionID).
		Update("revoked_at", time.Now()).Error
}
