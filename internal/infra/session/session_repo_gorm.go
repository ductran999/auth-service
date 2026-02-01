package session

import (
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
