package session_test

import (
	"testing"

	"auth-service/internal/apperrs"
	"auth-service/internal/biz/usecase/auth/credential"
	"auth-service/internal/biz/usecase/auth/session"
	"auth-service/internal/domain/accountmodel"
	"auth-service/test/fakes"
	mockbuilder "auth-service/test/mock-builder"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func NewAuthUseCaseUT(t *testing.T, builders *mockbuilder.BuilderContainer) session.AuthSessionUsecase {
	hasher := builders.HasherBuilder.GetInstance()
	accountRepo := builders.AccountRepoBuilder.GetInstance()

	return session.NewAuthSessionUsecase(
		credential.NewCredentialVerifier(hasher, accountRepo),
		builders.SessionStore.GetInstance(),
		builders.SessionRepoBuilder.GetInstance(),
	)
}

func TestLogin(t *testing.T) {
	loginInput := session.LoginInput{
		Email:    fakes.FakeAccount().Email,
		Password: fakes.FakeAccount().PasswordHash,
	}
	expectedAccount := fakes.FakeAccount()

	testTable := []struct {
		name        string
		setup       func(t *testing.T) session.AuthSessionUsecase
		loginInput  session.LoginInput
		expectedErr error
		expected    *accountmodel.Account
	}{
		{
			name: "verify account failed",
			loginInput: session.LoginInput{
				CurrentSessionID: uuid.Nil.String(),
				Email:            loginInput.Email,
				Password:         loginInput.Password,
			},
			setup: func(t *testing.T) session.AuthSessionUsecase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.AccountRepoBuilder.FindByEmailError()
				return NewAuthUseCaseUT(t, builders)
			},
			expectedErr: apperrs.ErrUnauthorized,
			expected:    nil,
		},
		{
			name:       "create session failed",
			loginInput: loginInput,
			setup: func(t *testing.T) session.AuthSessionUsecase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.AccountRepoBuilder.FindByEmailHasResult(t.Context(), fakes.FakeAccount().Email)
				builders.HasherBuilder.HashPasswordMatch(fakes.FakeAccount().PasswordHash)
				builders.SessionRepoBuilder.CreateSessionFailed()
				return NewAuthUseCaseUT(t, builders)
			},
			expectedErr: apperrs.ErrInternal,
			expected:    nil,
		},
		{
			name:       "set cache failed",
			loginInput: loginInput,
			setup: func(t *testing.T) session.AuthSessionUsecase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.AccountRepoBuilder.FindByEmailHasResult(t.Context(), fakes.FakeAccount().Email)
				builders.HasherBuilder.HashPasswordMatch(fakes.FakeAccount().PasswordHash)
				builders.SessionRepoBuilder.CreateSessionSuccess()
				builders.SessionStore.Set_Failed(t.Context())
				return NewAuthUseCaseUT(t, builders)
			},
			expectedErr: apperrs.ErrInternal,
			expected:    nil,
		},
		{
			name:       "login success",
			loginInput: loginInput,
			setup: func(t *testing.T) session.AuthSessionUsecase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.AccountRepoBuilder.FindByEmailHasResult(t.Context(), fakes.FakeAccount().Email)
				builders.HasherBuilder.HashPasswordMatch(fakes.FakeAccount().PasswordHash)
				builders.SessionRepoBuilder.CreateSessionSuccess()
				builders.SessionStore.Set_OK(t.Context())
				return NewAuthUseCaseUT(t, builders)
			},
			expectedErr: nil,
			expected:    expectedAccount,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			sut := tc.setup(t)
			ctx := t.Context()

			// Act
			session, err := sut.Login(ctx, tc.loginInput)

			// Assert
			assert.ErrorIs(t, err, tc.expectedErr)
			if tc.expected != nil {
				assert.Equal(t, tc.expected.ID, session.AccountID)
			} else {
				assert.Nil(t, session)
			}
		})
	}
}
