package mockbuilder

import (
	"auth-service/internal/domain/sessionmodel"
	"auth-service/test/fakes"
	"auth-service/test/mocks"
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
)

type mockSessionStore struct {
	inst *mocks.SessionStore
}

func newMockSessionStore(t *testing.T) *mockSessionStore {
	return &mockSessionStore{
		inst: mocks.NewSessionStore(t),
	}
}

func (b *mockSessionStore) GetInstance() *mocks.SessionStore {
	return b.inst
}

func (b *mockSessionStore) Get_Failed(ctx context.Context) {
	b.inst.EXPECT().
		Get(ctx, mock.AnythingOfType("string")).
		Return(nil, ErrInternalDB)
}

func (b *mockSessionStore) Get_OK(ctx context.Context) {
	session := sessionmodel.Session{
		AccountID: fakes.FakeAccount().ID,
	}

	b.inst.EXPECT().
		Get(ctx, mock.AnythingOfType("string")).
		Return(&session, nil)
}

func (b *mockSessionStore) Set_Failed(ctx context.Context) {
	b.inst.EXPECT().Save(ctx, mock.AnythingOfType("*sessionmodel.Session")).Return(ErrInternalDB)
}

func (b *mockSessionStore) Set_OK(ctx context.Context) {
	b.inst.EXPECT().Save(ctx, mock.AnythingOfType("*sessionmodel.Session")).Return(nil)
}
