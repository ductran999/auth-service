package port

import (
	"auth-service/internal/domain/sessionmodel"
	"context"
)

// SessionUsecase defines business logic for validating sessions in user-facing flows.
type SessionUsecase interface {
	// Validate checks if a session exists and is not expired.
	// It first looks in the cache, and if missing, checks persistent storage.
	Validate(ctx context.Context, sessionID string) (*sessionmodel.Session, error)
}
