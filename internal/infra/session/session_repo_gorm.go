package session

import (
	sessionUC "auth-service/internal/biz/usecase/session"
	"auth-service/internal/domain/sessionmodel"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

// sessionRepo implements the SessionRepository interface.
type sessionRepoImpl struct {
	db *gorm.DB
}

// NewSessionRepository returns a concrete implementation of SessionRepository.
func NewSessionRepository(db *gorm.DB) *sessionRepoImpl {
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

func (r *sessionRepoImpl) DeleteExpiredBefore(ctx context.Context, cutoff time.Time) error {
	query := `DELETE FROM sessions WHERE expires_at < ?`
	return r.db.WithContext(ctx).Exec(query, cutoff).Error
}

// MarkSessionsExpired sets the expiration timestamp for multiple sessions by their IDs.
func (r *sessionRepoImpl) MarkSessionsExpired(
	ctx context.Context,
	sessionIDs []string,
	expiresAt time.Time,
) error {
	if len(sessionIDs) == 0 {
		return nil // nothing to do
	}

	// Bulk update: set expires_at where session id is in sessionIDs
	err := r.db.WithContext(ctx).
		Model(&sessionmodel.Session{}).
		Where("id IN ?", sessionIDs).
		Update("expires_at", expiresAt).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *sessionRepoImpl) FindByID(ctx context.Context, sessionID string) (*sessionmodel.Session, error) {
	var session sessionmodel.Session

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

func (r *sessionRepoImpl) FindAllActiveSession(ctx context.Context) ([]sessionmodel.Session, error) {
	var activeSessions []sessionmodel.Session

	err := r.db.WithContext(ctx).
		Select("id").
		Where("expires_at IS NULL").
		Find(&activeSessions).Error

	if err != nil {
		return nil, err
	}

	return activeSessions, nil
}
