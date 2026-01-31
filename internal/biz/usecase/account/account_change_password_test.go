package account_test

import (
	"testing"

	"auth-service/internal/apperrs"
	"auth-service/internal/biz/usecase/account"
	"auth-service/test/fakes"
	mockbuilder "auth-service/test/mock-builder"

	"github.com/stretchr/testify/require"
)

func TestChangePassword(t *testing.T) {
	validInput := account.ChangePasswordInput{
		AccountID:   fakes.FakeAccount().ID.String(),
		OldPassword: mockbuilder.FakeOldPass,
		NewPassword: mockbuilder.FakeNewPass,
	}

	samePassInput := account.ChangePasswordInput{
		AccountID:   fakes.FakeAccount().ID.String(),
		OldPassword: fakes.FakeAccount().PasswordHash,
		NewPassword: fakes.FakeAccount().PasswordHash,
	}

	testTable := []struct {
		name        string
		setup       func(t *testing.T) account.AccountUsecase
		input       account.ChangePasswordInput
		expectedErr error
	}{
		{
			name: "failed find account by id",
			setup: func(t *testing.T) account.AccountUsecase {
				t.Helper()
				builders := mockbuilder.NewBuilderContainer(t)
				builders.AccountRepoBuilder.FindByIDFailed()
				return NewAccountUseCaseUT(t, builders)
			},
			input:       validInput,
			expectedErr: apperrs.ErrInternal,
		},
		{
			name: "not found account by id",
			setup: func(t *testing.T) account.AccountUsecase {
				t.Helper()
				builders := mockbuilder.NewBuilderContainer(t)
				builders.AccountRepoBuilder.FindByID_NoResult()
				return NewAccountUseCaseUT(t, builders)
			},
			input:       validInput,
			expectedErr: apperrs.ErrNotFound,
		},
		{
			name: "failed to compare old password to hash",
			setup: func(t *testing.T) account.AccountUsecase {
				t.Helper()
				builders := mockbuilder.NewBuilderContainer(t)
				builders.AccountRepoBuilder.FindByIDSuccess(t.Context(), fakes.FakeAccount().ID.String())
				builders.HasherBuilder.CompareHashPasswordGotError()
				return NewAccountUseCaseUT(t, builders)
			},
			input:       validInput,
			expectedErr: apperrs.ErrInternal,
		},
		{
			name: "old password not match",
			setup: func(t *testing.T) account.AccountUsecase {
				t.Helper()
				builders := mockbuilder.NewBuilderContainer(t)
				builders.AccountRepoBuilder.FindByIDSuccess(t.Context(), fakes.FakeAccount().ID.String())
				builders.HasherBuilder.HashPasswordNotMatch()
				return NewAccountUseCaseUT(t, builders)
			},
			input:       validInput,
			expectedErr: apperrs.ErrUnauthorized,
		},
		{
			name:  "failed to hashing password",
			input: validInput,
			setup: func(t *testing.T) account.AccountUsecase {
				t.Helper()
				builders := mockbuilder.NewBuilderContainer(t)
				builders.AccountRepoBuilder.FindByIDSuccess(t.Context(), fakes.FakeAccount().ID.String())
				builders.HasherBuilder.HashPasswordMatch(fakes.FakeAccount().PasswordHash)
				builders.HasherBuilder.HashingPasswordFailed()
				return NewAccountUseCaseUT(t, builders)
			},
			expectedErr: apperrs.ErrInternal,
		},
		{
			name: "failed to update new password",
			setup: func(t *testing.T) account.AccountUsecase {
				t.Helper()
				builders := mockbuilder.NewBuilderContainer(t)
				builders.AccountRepoBuilder.FindByIDSuccess(t.Context(), fakes.FakeAccount().ID.String())
				builders.HasherBuilder.HashPasswordMatch(fakes.FakeAccount().PasswordHash)
				builders.HasherBuilder.HashingPasswordSuccess()
				builders.AccountRepoBuilder.UpdatePasswordHashFailed()
				return NewAccountUseCaseUT(t, builders)
			},
			input:       validInput,
			expectedErr: apperrs.ErrInternal,
		},
		{
			name: "new password must same ass the old one",
			setup: func(t *testing.T) account.AccountUsecase {
				t.Helper()
				builders := mockbuilder.NewBuilderContainer(t)
				builders.AccountRepoBuilder.FindByIDSuccess(t.Context(), fakes.FakeAccount().ID.String())
				builders.HasherBuilder.HashPasswordMatch(fakes.FakeAccount().PasswordHash)
				return NewAccountUseCaseUT(t, builders)
			},
			input:       samePassInput,
			expectedErr: apperrs.ErrInvalidInput,
		},
		{
			name: "change password success",
			setup: func(t *testing.T) account.AccountUsecase {
				t.Helper()
				builders := mockbuilder.NewBuilderContainer(t)
				builders.AccountRepoBuilder.FindByIDSuccess(t.Context(), fakes.FakeAccount().ID.String())
				builders.HasherBuilder.HashPasswordMatch(fakes.FakeAccount().PasswordHash)
				builders.HasherBuilder.HashingPasswordSuccess()
				builders.AccountRepoBuilder.UpdatePasswordHashSuccess()
				return NewAccountUseCaseUT(t, builders)
			},
			input:       validInput,
			expectedErr: nil,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			sut := tc.setup(t)

			err := sut.ChangePassword(t.Context(), tc.input)

			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
