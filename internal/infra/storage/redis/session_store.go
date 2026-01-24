package redis

type sessionStore struct{}

func NewSessionStore() *sessionStore {
	return &sessionStore{}
}
