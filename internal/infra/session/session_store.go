package session

import (
	"auth-service/internal/biz/ports/auth"
	"auth-service/internal/domain/sessionmodel"
	"context"
	"time"

	"github.com/DucTran999/cachekit"
)

const (
	SessionKeyPrefix = "auth:session:"
	SessionTTL       = time.Hour
)

type sessionStore struct {
	redisClient cachekit.RemoteCache
}

func NewSessionStore(cl cachekit.RemoteCache) auth.SessionStore {
	return &sessionStore{redisClient: cl}
}

func (s *sessionStore) Get(ctx context.Context, id string) (*sessionmodel.Session, error) {
	key := s.cacheKey(id)

	session := &sessionmodel.Session{}
	if err := s.redisClient.GetInto(ctx, key, session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *sessionStore) Save(ctx context.Context, session *sessionmodel.Session) error {
	key := s.cacheKey(session.ID.String())
	return s.redisClient.Set(ctx, key, session, SessionTTL)
}

func (s *sessionStore) Refresh(ctx context.Context, id string) (*sessionmodel.Session, error) {
	return nil, nil
}

func (s *sessionStore) cacheKey(sessionID string) string {
	return SessionKeyPrefix + sessionID
}
