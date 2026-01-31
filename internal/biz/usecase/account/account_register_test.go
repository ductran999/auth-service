package account_test

import (
	"testing"

	"auth-service/internal/apperrs"
	"auth-service/internal/biz/usecase/account"
	"auth-service/internal/domain/accountmodel"
	"auth-service/test/fakes"
	mockbuilder "auth-service/test/mock-builder"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NewAccountUseCaseUT(t *testing.T, builders *mockbuilder.BuilderContainer) account.AccountUsecase {
	return account.NewAccountUseCase(
		builders.HasherBuilder.GetInstance(),
		builders.AccountRepoBuilder.GetInstance(),
	)
}

func TestRegisterAccount(t *testing.T) {
	userSample := account.RegisterInput{
		Email:    fakes.FakeAccount().Email,
		Password: "some-password",
	}

	testTable := []struct {
		name        string
		setup       func(t *testing.T) account.AccountUsecase
		accountInfo account.RegisterInput
		expectedErr error
		expected    *accountmodel.Account
	}{
		{
			name: "failed to find email in db",
			setup: func(t *testing.T) account.AccountUsecase {
				t.Helper()
				b := mockbuilder.NewBuilderContainer(t)
				b.AccountRepoBuilder.FindByEmailError()
				return NewAccountUseCaseUT(t, b)
			},
			accountInfo: userSample,
			expectedErr: apperrs.ErrInternal,
			expected:    nil,
		},
		{
			name: "failed caused email already taken",
			setup: func(t *testing.T) account.AccountUsecase {
				t.Helper()
				b := mockbuilder.NewBuilderContainer(t)
				b.AccountRepoBuilder.FindByEmailHasResult(t.Context(), fakes.FakeAccount().Email)
				return NewAccountUseCaseUT(t, b)
			},
			accountInfo: userSample,
			expectedErr: apperrs.ErrConflict,
			expected:    nil,
		},
		{
			name: "failed when hash password",
			setup: func(t *testing.T) account.AccountUsecase {
				t.Helper()
				b := mockbuilder.NewBuilderContainer(t)
				b.AccountRepoBuilder.FindByEmailNoResult()
				b.HasherBuilder.HashingPasswordFailed()
				return NewAccountUseCaseUT(t, b)
			},
			accountInfo: userSample,
			expectedErr: apperrs.ErrInternal,
			expected:    nil,
		},
		{
			name: "failed when persist to db",
			setup: func(t *testing.T) account.AccountUsecase {
				t.Helper()
				b := mockbuilder.NewBuilderContainer(t)
				b.AccountRepoBuilder.FindByEmailNoResult()
				b.HasherBuilder.HashingPasswordSuccess()
				b.AccountRepoBuilder.CreateAccountError()
				return NewAccountUseCaseUT(t, b)
			},
			accountInfo: userSample,
			expectedErr: apperrs.ErrInternal,
			expected:    nil,
		},
		{
			name: "register success",
			setup: func(t *testing.T) account.AccountUsecase {
				t.Helper()
				b := mockbuilder.NewBuilderContainer(t)
				b.AccountRepoBuilder.FindByEmailNoResult()
				b.HasherBuilder.HashingPasswordSuccess()
				b.AccountRepoBuilder.CreateAccountSuccess()
				return NewAccountUseCaseUT(t, b)
			},
			accountInfo: userSample,
			expectedErr: nil,
			expected:    fakes.FakeAccount(),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			sut := tc.setup(t)
			user, err := sut.Register(t.Context(), tc.accountInfo)

			require.ErrorIs(t, err, tc.expectedErr)
			if tc.expected != nil {
				assert.Equal(t, tc.expected.Email, user.Email)
			} else {
				assert.Nil(t, user)
			}
		})
	}
}
