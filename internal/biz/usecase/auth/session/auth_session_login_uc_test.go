package session_test

import (
	"testing"

	"auth-service/internal/apperrs"
	"auth-service/internal/biz/usecase/auth/credential"
	"auth-service/internal/biz/usecase/auth/session"
	"auth-service/internal/domain/accountmodel"
	mockbuilder "auth-service/test/mock-builder"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func NewAuthUseCaseUT(t *testing.T, builders *mockbuilder.BuilderContainer) session.AuthSessionUsecase {
	hasher := builders.HasherBuilder.GetInstance()
	accountRepo := builders.AccountRepoBuilder.GetInstance()

	return session.NewAuthSessionUsecase(
		credential.NewCredentialVerifier(hasher, accountRepo),
		nil,
		builders.SessionRepoBuilder.GetInstance(),
	)
}

func TestLogin(t *testing.T) {
	type testCase struct {
		name        string
		setup       func(t *testing.T) session.AuthSessionUsecase
		loginInput  session.LoginInput
		expectedErr error
		expected    *accountmodel.Account
	}

	loginInput := session.LoginInput{
		Email:    mockbuilder.FakeEmail,
		Password: mockbuilder.FakeOldPass,
	}
	expectedAccount := &accountmodel.Account{
		ID:       mockbuilder.FakeAccountID,
		Email:    mockbuilder.FakeEmail,
		IsActive: true,
	}

	testTable := []testCase{
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
				builders.AccountRepoBuilder.FindByEmailHasResult()
				builders.HasherBuilder.HashPasswordMatch()
				builders.SessionRepoBuilder.CreateSessionFailed()
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
				builders.AccountRepoBuilder.FindByEmailHasResult()
				builders.HasherBuilder.HashPasswordMatch()
				builders.SessionRepoBuilder.CreateSessionSuccess()
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
