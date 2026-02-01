package account_test

import (
	"testing"

	"auth-service/internal/apperrs"
	"auth-service/internal/model"
	"auth-service/internal/usecase/account"
	"auth-service/internal/usecase/dto"
	"auth-service/internal/usecase/port"
	mockbuilder "auth-service/test/mock-builder"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NewAccountUseCaseUT(t *testing.T, builders *mockbuilder.BuilderContainer) port.AccountUsecase {
	return account.NewAccountUseCase(
		builders.HasherBuilder.GetInstance(),
		builders.AccountRepoBuilder.GetInstance(),
	)
}

func TestRegisterAccount(t *testing.T) {
	type testCase struct {
		name        string
		setup       func(t *testing.T) port.AccountUsecase
		accountInfo dto.RegisterInput
		expectedErr error
		expected    *model.Account
	}

	userSample := dto.RegisterInput{
		Email:    mockbuilder.FakeEmail,
		Password: "abc1234!",
	}

	testTable := []testCase{
		{
			name: "failed to find email in db",
			setup: func(t *testing.T) port.AccountUsecase {
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
			setup: func(t *testing.T) port.AccountUsecase {
				t.Helper()
				b := mockbuilder.NewBuilderContainer(t)
				b.AccountRepoBuilder.FindByEmailHasResult()
				return NewAccountUseCaseUT(t, b)
			},
			accountInfo: userSample,
			expectedErr: apperrs.ErrConflict,
			expected:    nil,
		},
		{
			name: "failed when hash password",
			setup: func(t *testing.T) port.AccountUsecase {
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
			setup: func(t *testing.T) port.AccountUsecase {
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
			setup: func(t *testing.T) port.AccountUsecase {
				t.Helper()
				b := mockbuilder.NewBuilderContainer(t)
				b.AccountRepoBuilder.FindByEmailNoResult()
				b.HasherBuilder.HashingPasswordSuccess()
				b.AccountRepoBuilder.CreateAccountSuccess()
				return NewAccountUseCaseUT(t, b)
			},
			accountInfo: userSample,
			expectedErr: nil,
			expected: &model.Account{
				Email: "daniel@example.com",
			},
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
