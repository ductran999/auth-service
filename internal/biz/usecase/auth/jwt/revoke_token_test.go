package jwt_test

import (
	"auth-service/internal/apperrs"
	"auth-service/internal/biz/usecase/auth/jwt"
	"auth-service/test/fakes"
	mockbuilder "auth-service/test/mock-builder"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRevokeToken(t *testing.T) {
	refreshToken := fakes.FakeTokenPairs().RefreshToken

	testcases := []struct {
		name         string
		refreshToken string
		setup        func(t *testing.T) jwt.AuthJWTUsecase
		expectedErr  error
	}{
		{
			name: "missing refresh token",
			setup: func(t *testing.T) jwt.AuthJWTUsecase {
				t.Helper()
				builder := mockbuilder.NewBuilderContainer(t)
				return NewAuthJWTUsecaseUT(t, builder)
			},
			expectedErr: apperrs.ErrUnauthorized,
		},
		{
			name:         "invalid refresh token",
			refreshToken: refreshToken,
			setup: func(t *testing.T) jwt.AuthJWTUsecase {
				t.Helper()
				builder := mockbuilder.NewBuilderContainer(t)
				builder.TokenService.VerifyRefreshToken_Failed(t.Context())
				return NewAuthJWTUsecaseUT(t, builder)
			},
			expectedErr: apperrs.ErrUnauthorized,
		},
		{
			name:         "revoke failed",
			refreshToken: refreshToken,
			setup: func(t *testing.T) jwt.AuthJWTUsecase {
				t.Helper()
				builder := mockbuilder.NewBuilderContainer(t)
				builder.TokenService.VerifyRefreshToken_Success(t.Context())
				builder.TokenStore.Revoke_Failed(t.Context())
				return NewAuthJWTUsecaseUT(t, builder)
			},
			expectedErr: apperrs.ErrInternal,
		},
		{
			name:         "revoke success",
			refreshToken: refreshToken,
			setup: func(t *testing.T) jwt.AuthJWTUsecase {
				t.Helper()
				builder := mockbuilder.NewBuilderContainer(t)
				builder.TokenService.VerifyRefreshToken_Success(t.Context())
				builder.TokenStore.Revoke_OK(t.Context())
				return NewAuthJWTUsecaseUT(t, builder)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			sut := tc.setup(t)

			// Act
			err := sut.RevokeRefreshToken(t.Context(), tc.refreshToken)

			// Assert
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
