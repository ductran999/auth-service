package credential_test

import (
	"auth-service/internal/biz/usecase/auth/credential"
	"auth-service/internal/domain/accountmodel"
	mockbuilder "auth-service/test/mock-builder"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVerify_TableDriven(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		setup       func(t *testing.T) credential.CredentialVerifier
		email       string
		password    string
		expectedErr error
		expected    *accountmodel.Account
	}{
		{
			name:     "invalid credentials",
			email:    "test@example.com",
			password: "password",
			setup: func(t *testing.T) credential.CredentialVerifier {
				t.Helper()
				builders := mockbuilder.NewBuilderContainer(t)
				builders.AccountRepoBuilder.FindByEmailNoResult()

				return *credential.NewCredentialVerifier(
					builders.HasherBuilder.GetInstance(),
					builders.AccountRepoBuilder.GetInstance(),
				)
			},
			expectedErr: credential.ErrInvalidCredentials,
		},
		{
			name:     "disable account",
			email:    "test@example.com",
			password: "password",
			setup: func(t *testing.T) credential.CredentialVerifier {
				t.Helper()
				builders := mockbuilder.NewBuilderContainer(t)
				builders.AccountRepoBuilder.FindByEmailAccountInactive()

				return *credential.NewCredentialVerifier(
					builders.HasherBuilder.GetInstance(),
					builders.AccountRepoBuilder.GetInstance(),
				)
			},
			expectedErr: credential.ErrAccountDisabled,
		},
		{
			name:     "compare hash got error",
			email:    "test@example.com",
			password: "password",
			setup: func(t *testing.T) credential.CredentialVerifier {
				t.Helper()
				builders := mockbuilder.NewBuilderContainer(t)
				builders.AccountRepoBuilder.FindByEmailHasResult()
				builders.HasherBuilder.CompareHashPasswordGotError()
				return *credential.NewCredentialVerifier(
					builders.HasherBuilder.GetInstance(),
					builders.AccountRepoBuilder.GetInstance(),
				)
			},
			expectedErr: mockbuilder.ErrCompareHashPassword,
		},
		{
			name:     "password not match",
			email:    "test@example.com",
			password: "password",
			setup: func(t *testing.T) credential.CredentialVerifier {
				t.Helper()
				builders := mockbuilder.NewBuilderContainer(t)
				builders.AccountRepoBuilder.FindByEmailHasResult()
				builders.HasherBuilder.HashPasswordNotMatch()
				return *credential.NewCredentialVerifier(
					builders.HasherBuilder.GetInstance(),
					builders.AccountRepoBuilder.GetInstance(),
				)
			},
			expectedErr: credential.ErrInvalidCredentials,
		},
		{
			name:     "verify success",
			email:    "test@example.com",
			password: "password",
			setup: func(t *testing.T) credential.CredentialVerifier {
				t.Helper()
				builders := mockbuilder.NewBuilderContainer(t)
				builders.AccountRepoBuilder.FindByEmailHasResult()
				builders.HasherBuilder.HashPasswordMatch()
				return *credential.NewCredentialVerifier(
					builders.HasherBuilder.GetInstance(),
					builders.AccountRepoBuilder.GetInstance(),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			verifier := tt.setup(t)
			_, err := verifier.Verify(t.Context(), tt.email, tt.password)
			require.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
