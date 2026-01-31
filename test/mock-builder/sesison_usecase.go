package mockbuilder

import (
	"errors"
	"testing"

	"auth-service/internal/biz/usecase/session"
	"auth-service/internal/domain/sessionmodel"
	"auth-service/test/fakes"
	"auth-service/test/mocks"

	"github.com/stretchr/testify/mock"
)

var (
	ErrSessionValidate = errors.New("failed to validate session")
)

type mockSessionUsecase struct {
	inst *mocks.SessionUsecase
}

func newMockSessionUsecase(t *testing.T) *mockSessionUsecase {
	return &mockSessionUsecase{
		inst: mocks.NewSessionUsecase(t),
	}
}

func (m *mockSessionUsecase) GetInstance() *mocks.SessionUsecase {
	return m.inst
}

func (m *mockSessionUsecase) ValidateError() {
	m.inst.EXPECT().
		Validate(mock.Anything, mock.Anything).
		Return(nil, ErrSessionValidate)
}

func (m *mockSessionUsecase) ValidateInvalidSession() {
	m.inst.EXPECT().
		Validate(mock.Anything, mock.Anything).
		Return(nil, session.ErrInvalidSessionID)
}

func (m *mockSessionUsecase) ValidateSessionSuccess() {
	m.inst.EXPECT().
		Validate(mock.Anything, mock.Anything).
		Return(&sessionmodel.Session{
			ID:        FakeSessionID,
			AccountID: fakes.FakeAccount().ID,
		}, nil)
}
