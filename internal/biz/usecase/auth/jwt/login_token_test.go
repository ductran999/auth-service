package jwt_test

import (
	"auth-service/internal/apperrs"
	"auth-service/internal/biz/usecase/auth/credential"
	"auth-service/internal/biz/usecase/auth/jwt"
	"auth-service/test/fakes"
	mockbuilder "auth-service/test/mock-builder"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NewAuthJWTUsecaseUT(t *testing.T, builders *mockbuilder.BuilderContainer) jwt.AuthJWTUsecase {
	hasher := builders.HasherBuilder.GetInstance()
	accountRepo := builders.AccountRepoBuilder.GetInstance()

	return jwt.NewAuthJWTUsecase(
		builders.TokenService.GetInstance(),
		builders.TokenStore.GetInstance(),
		credential.NewCredentialVerifier(hasher, accountRepo),
	)
}

func TestLoginWithJWT(t *testing.T) {
	testcases := []struct {
		name        string
		setup       func(t *testing.T) jwt.AuthJWTUsecase
		expectedErr error
	}{
		{
			name: "authenticate failed",
			setup: func(t *testing.T) jwt.AuthJWTUsecase {
				t.Helper()
				builder := mockbuilder.NewBuilderContainer(t)
				builder.AccountRepoBuilder.FindByEmailNoResult()
				return NewAuthJWTUsecaseUT(t, builder)
			},
			expectedErr: apperrs.ErrUnauthorized,
		},
		{
			name: "sign access token failed",
			setup: func(t *testing.T) jwt.AuthJWTUsecase {
				t.Helper()
				builder := mockbuilder.NewBuilderContainer(t)
				builder.AccountRepoBuilder.FindByEmailHasResult(t.Context(), fakes.FakeAccount().Email)
				builder.HasherBuilder.HashPasswordMatch(fakes.FakeAccount().PasswordHash)
				builder.TokenService.SignPairs_Failed(fakes.FakeAccount())
				return NewAuthJWTUsecaseUT(t, builder)
			},
			expectedErr: apperrs.ErrInternal,
		},
		// {
		// 	name: "sign refresh token failed",
		// 	setup: func(t *testing.T) port.AuthJWTUsecase {
		// 		t.Helper()
		// 		builder := mockbuilder.NewBuilderContainer(t)
		// 		builder.AccountVerifier.VerifySuccess()
		// 		builder.TokenSigner.SignAccessSuccessAndSignRefreshFailed()
		// 		return NewAuthJWTUsecaseUT(t, builder)
		// 	},
		// 	expectedErr: mockbuilder.ErrSigningToken,
		// },
		{
			name: "save session failed",
			setup: func(t *testing.T) jwt.AuthJWTUsecase {
				t.Helper()
				account := fakes.FakeAccount()
				builder := mockbuilder.NewBuilderContainer(t)
				builder.AccountRepoBuilder.FindByEmailHasResult(t.Context(), account.Email)
				builder.HasherBuilder.HashPasswordMatch(account.PasswordHash)
				builder.TokenService.SignPairs_Success(account)
				builder.TokenStore.Save_Failed(t.Context())
				return NewAuthJWTUsecaseUT(t, builder)
			},
			expectedErr: apperrs.ErrInternal,
		},
		{
			name: "login success",
			setup: func(t *testing.T) jwt.AuthJWTUsecase {
				t.Helper()
				account := fakes.FakeAccount()
				builder := mockbuilder.NewBuilderContainer(t)
				builder.AccountRepoBuilder.FindByEmailHasResult(t.Context(), account.Email)
				builder.HasherBuilder.HashPasswordMatch(account.PasswordHash)
				builder.TokenService.SignPairs_Success(account)
				builder.TokenStore.Save_Success(t.Context())
				return NewAuthJWTUsecaseUT(t, builder)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			input := jwt.LoginJWTInput{
				Email:    fakes.FakeAccount().Email,
				Password: fakes.FakeAccount().PasswordHash,
			}
			sut := tc.setup(t)

			// Act
			tokens, err := sut.Login(t.Context(), input)

			// Assert
			require.ErrorIs(t, err, tc.expectedErr)
			if tc.expectedErr == nil {
				require.NotEmpty(t, tokens.AccessToken)
				assert.NotEmpty(t, tokens.RefreshToken)
			}
		})
	}
}
