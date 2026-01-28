package mockbuilder

import (
	"context"
	"errors"
	"testing"

	"auth-service/internal/biz/usecase/account"
	"auth-service/internal/domain/accountmodel"
	"auth-service/test/mocks"

	"github.com/stretchr/testify/mock"
)

var (
	FakeEmail   = "daniel@example.com"
	FakeOldPass = "0ldP@ssW0rd"
	FakeNewPass = "N3wP@ssW0rd"

	ErrFindAccountByEmail = errors.New("find email unexpected error")
	ErrFindAccountByID    = errors.New("find id unexpected error")
	ErrCreateAccount      = errors.New("unexpected error create new account")
	ErrUpdateHashPassword = errors.New("unexpected error when update hash password")
)

type mockAccountRepoBuilder struct {
	inst *mocks.AccountRepo
}

func newMockAccountRepoBuilder(t *testing.T) *mockAccountRepoBuilder {
	return &mockAccountRepoBuilder{
		inst: mocks.NewAccountRepo(t),
	}
}

func (b *mockAccountRepoBuilder) GetInstance() *mocks.AccountRepo {
	return b.inst
}

func (b *mockAccountRepoBuilder) CreateAccountError() {
	b.inst.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(ErrCreateAccount)
}

func (b *mockAccountRepoBuilder) CreateAccountSuccess() {
	mockAccount := &accountmodel.Account{
		ID:       FakeAccountID,
		Email:    FakeEmail,
		IsActive: true,
	}
	b.inst.EXPECT().
		Create(mock.Anything, mock.Anything).Run(
		func(ctx context.Context, account *accountmodel.Account) {
			*account = *mockAccount
		},
	).Return(nil)
}

func (b *mockAccountRepoBuilder) FindByEmailError() {
	b.inst.EXPECT().
		FindByEmail(mock.Anything, mock.Anything).
		Return(nil, ErrFindAccountByEmail)
}

func (b *mockAccountRepoBuilder) FindByEmailHasResult() {
	activeAccount := &accountmodel.Account{
		ID:       FakeAccountID,
		Email:    FakeEmail,
		IsActive: true,
	}

	b.inst.EXPECT().
		FindByEmail(mock.Anything, mock.Anything).
		Return(activeAccount, nil)
}

func (b *mockAccountRepoBuilder) FindByEmailAccountInactive() {
	mockAccount := &accountmodel.Account{
		ID:       FakeAccountID,
		Email:    FakeEmail,
		IsActive: false,
	}

	b.inst.EXPECT().
		FindByEmail(mock.Anything, mock.Anything).
		Return(mockAccount, nil)
}

func (b *mockAccountRepoBuilder) FindByEmailNoResult() {
	b.inst.EXPECT().
		FindByEmail(mock.Anything, mock.Anything).
		Return(nil, account.ErrAccountNotFound)
}

func (b *mockAccountRepoBuilder) FindByIDFailed() {
	b.inst.EXPECT().
		FindByID(mock.Anything, mock.Anything).
		Return(nil, ErrFindAccountByID)
}

func (b *mockAccountRepoBuilder) FindByIDSuccess() {
	mockAccount := &accountmodel.Account{
		ID:           FakeAccountID,
		Email:        FakeEmail,
		PasswordHash: FakeOldPass,
		IsActive:     true,
	}

	b.inst.EXPECT().
		FindByID(mock.Anything, mock.Anything).
		Return(mockAccount, nil)
}

func (b *mockAccountRepoBuilder) UpdatePasswordHashFailed() {
	b.inst.EXPECT().
		UpdatePasswordHash(mock.Anything, mock.Anything, mock.Anything).
		Return(ErrUpdateHashPassword)
}

func (b *mockAccountRepoBuilder) UpdatePasswordHashSuccess() {
	b.inst.EXPECT().
		UpdatePasswordHash(mock.Anything,
			mock.MatchedBy(func(id string) bool {
				return id == FakeAccountID.String()
			}),
			mock.Anything).
		Return(nil)
}
