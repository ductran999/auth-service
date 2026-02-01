package repositories

import (
	"auth-service/internal/domain/sessionmodel"
	"context"
)

// SessionRepository defines the methods required to manage session data in the persistence layer.
type SessionRepository interface {
	// Create stores a new session in the database.
	Create(ctx context.Context, session *sessionmodel.Session) error

	// Revoke marks a session as revoked by updating its expiration timestamp.
	Revoke(ctx context.Context, sessionID string) error
}
