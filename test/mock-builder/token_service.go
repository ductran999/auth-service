package mockbuilder

import (
	"auth-service/test/mocks"
	"testing"
)

type mockTokenService struct {
	inst *mocks.TokenService
}

func (b *mockTokenService) GetInstance() *mocks.TokenService {
	return b.inst
}

func newMockTokenService(t *testing.T) *mockTokenService {
	return &mockTokenService{
		inst: mocks.NewTokenService(t),
	}
}
