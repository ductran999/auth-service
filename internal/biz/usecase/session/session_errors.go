package session

import "errors"

var (
	ErrSessionNotFound  = errors.New("session not found or expired")
	ErrInvalidSessionID = errors.New("invalid session id")
)
