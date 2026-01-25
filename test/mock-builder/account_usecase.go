package mockbuilder

import (
	"errors"
	"testing"

	"auth-service/internal/apperrs"
	"auth-service/internal/model"
	"auth-service/internal/usecase/account"
	"auth-service/test/mocks"

	"github.com/stretchr/testify/mock"
)

var (
	ErrRegisterAccount = errors.New("failed to register account")
	ErrChangePassword  = errors.New("failed to change password")
)

type mockAccountUsecase struct {
	inst *mocks.AccountUsecase
}

func newMockAccountUsecase(t *testing.T) *mockAccountUsecase {
	return &mockAccountUsecase{
		inst: mocks.NewAccountUsecase(t),
	}
}

func (m *mockAccountUsecase) GetInstance() *mocks.AccountUsecase {
	return m.inst
}

func (m *mockAccountUsecase) RegisterError() {
	m.inst.EXPECT().
		Register(mock.Anything, mock.Anything).
		Return(nil, ErrRegisterAccount)
}

func (m *mockAccountUsecase) RegisterConflictEmail() {
	m.inst.EXPECT().
		Register(mock.Anything, mock.Anything).
		Return(nil, apperrs.Conflict(account.ErrEmailExisted))
}

func (m *mockAccountUsecase) RegisterSuccess() {
	m.inst.EXPECT().
		Register(mock.Anything, mock.Anything).
		Return(&model.Account{
			ID: FakeAccountID,
		}, nil)
}

func (m *mockAccountUsecase) ChangePasswordGotErrorSamePass() {
	m.inst.EXPECT().
		ChangePassword(mock.Anything, mock.Anything).
		Return(account.ErrNewPasswordMustChanged)
}

func (m *mockAccountUsecase) ChangePasswordSuccess() {
	m.inst.EXPECT().
		ChangePassword(mock.Anything, mock.Anything).
		Return(nil)
}

func (m *mockAccountUsecase) ChangePassErrGotWrongCredentials() {
	m.inst.EXPECT().
		ChangePassword(mock.Anything, mock.Anything).
		Return(account.ErrPasswordMismatch)
}

func (m *mockAccountUsecase) ChangePassErrGotErrorDB() {
	m.inst.EXPECT().
		ChangePassword(mock.Anything, mock.Anything).
		Return(ErrChangePassword)
}
