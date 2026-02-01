package session

import (
	"auth-service/internal/model"
	"context"
	"time"
)

// Create inserts a new session record into the database.
func (r *sessionRepoImpl) Create(ctx context.Context, session *model.Session) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *sessionRepoImpl) Revoke(ctx context.Context, sessionID string) error {
	return r.db.WithContext(ctx).
		Model(&model.Session{}).
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
		Model(&model.Session{}).
		Where("id IN ?", sessionIDs).
		Update("expires_at", expiresAt).Error
	if err != nil {
		return err
	}

	return nil
}
