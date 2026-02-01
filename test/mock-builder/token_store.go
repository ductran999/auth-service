package mockbuilder

import (
	"auth-service/test/mocks"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
)

var (
	ErrInternalDB = errors.New("internal error db")
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

func (b *mockTokenStore) Save_Failed(ctx context.Context) {
	b.inst.EXPECT().Save(ctx, mock.Anything).Return(ErrInternalDB)
}

func (b *mockTokenStore) Save_Success(ctx context.Context) {
	b.inst.EXPECT().Save(ctx, mock.Anything).Return(nil)
}

func (b *mockTokenStore) Exists_Failed(ctx context.Context) {
	b.inst.EXPECT().Exists(ctx, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(false, ErrInternalDB)
}

func (b *mockTokenStore) Exists_NotExisted(ctx context.Context) {
	b.inst.EXPECT().Exists(ctx, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(false, nil)
}

func (b *mockTokenStore) Exists_OK(ctx context.Context) {
	b.inst.EXPECT().Exists(ctx, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(true, nil)
}

func (b *mockTokenStore) Revoke_OK(ctx context.Context) {
	b.inst.EXPECT().
		Revoke(ctx, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(nil)
}

func (b *mockTokenStore) Revoke_Failed(ctx context.Context) {
	b.inst.EXPECT().
		Revoke(ctx, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(ErrInternalDB)
}
