package mockbuilder

import (
	"auth-service/test/mocks"
	"testing"
)

type mockTokenStore struct {
	inst *mocks.TokenStore
}

func (b *mockTokenStore) GetInstance() *mocks.TokenStore {
	return b.inst
}

func newMockTokenStore(t *testing.T) *mockTokenStore {
	return &mockTokenStore{
		inst: mocks.NewTokenStore(t),
	}
}
