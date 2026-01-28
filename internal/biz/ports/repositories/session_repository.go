package repositories

import (
	"auth-service/internal/domain/sessionmodel"
	"context"
	"time"
)

// SessionRepository defines the methods required to manage session data in the persistence layer.
type SessionRepository interface {
	// Create stores a new session in the database.
	Create(ctx context.Context, session *sessionmodel.Session) error

	// DeleteExpiredBefore permanently deletes sessions that expired before the given cutoff time.
	DeleteExpiredBefore(ctx context.Context, cutoff time.Time) error

	// FindAllActiveSession retrieves all currently active sessions.
	FindAllActiveSession(ctx context.Context) ([]sessionmodel.Session, error)

	// FindByID retrieves a session by its session ID.
	// Returns nil if the session is not found.
	FindByID(ctx context.Context, sessionID string) (*sessionmodel.Session, error)

	// Revoke marks a session as revoked by updating its expiration timestamp.
	Revoke(ctx context.Context, sessionID string) error

	// MarkSessionsExpired sets the expiration timestamp for multiple sessions by their IDs.
	MarkSessionsExpired(ctx context.Context, sessionIDs []string, expiresAt time.Time) error
}
