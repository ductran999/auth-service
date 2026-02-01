package mockbuilder

import (
	"context"
	"errors"
	"testing"

	"auth-service/internal/domain/sessionmodel"
	"auth-service/test/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

var (
	FakeSessionID = uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	ErrCreateSession        = errors.New("unexpected create session error")
	ErrFindSessionByID      = errors.New("unexpected find session error")
	ErrUpdateSessionExpires = errors.New("unexpected update expires error")
	ErrDeleteExpiredBefore  = errors.New("unexpected error when delete session from db")
	ErrFindActiveSession    = errors.New("unexpected error while querying list active session")
	ErrMarkSessionsExpired  = errors.New("unexpected error while update sessions to expired")
)

type mockSessionRepoBuilder struct {
	inst *mocks.SessionRepository
}

func newMockSessionRepoBuilder(t *testing.T) *mockSessionRepoBuilder {
	return &mockSessionRepoBuilder{
		inst: mocks.NewSessionRepository(t),
	}
}

func (b *mockSessionRepoBuilder) GetInstance() *mocks.SessionRepository {
	return b.inst
}

func (blr *mockSessionRepoBuilder) RevokeFailed() {
	blr.inst.EXPECT().Revoke(mock.Anything, mock.Anything).Return(ErrUpdateSessionExpires)
}

func (blr *mockSessionRepoBuilder) RevokeSuccess() {
	blr.inst.EXPECT().Revoke(mock.Anything, mock.Anything).Return(nil)
}

func (blr *mockSessionRepoBuilder) CreateSessionSuccess() {
	blr.inst.EXPECT().
		Create(mock.Anything, mock.AnythingOfType("*sessionmodel.Session")).
		Run(func(ctx context.Context, session *sessionmodel.Session) {
			session.ID = FakeSessionID
		}).
		Return(nil)
}

func (blr *mockSessionRepoBuilder) CreateSessionFailed() {
	blr.inst.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(ErrCreateSession)
}
