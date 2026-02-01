package session_test

import (
	"auth-service/internal/apperrs"
	"auth-service/internal/biz/usecase/auth/session"
	"auth-service/test/fakes"
	mockbuilder "auth-service/test/mock-builder"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateSession(t *testing.T) {
	testTable := []struct {
		name      string
		sessionID string
		setup     func(t *testing.T) session.AuthSessionUsecase
		expectErr error
	}{
		{
			name:      "invalid session",
			sessionID: "",
			setup: func(t *testing.T) session.AuthSessionUsecase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.SessionStore.Get_Failed(t.Context())
				return NewAuthUseCaseUT(t, builders)
			},
			expectErr: apperrs.ErrUnauthorized,
		},
		{
			name:      "valid session",
			sessionID: fakes.Session().ID.String(),
			setup: func(t *testing.T) session.AuthSessionUsecase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.SessionStore.Get_OK(t.Context())
				return NewAuthUseCaseUT(t, builders)
			},
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			sut := tc.setup(t)

			_, err := sut.ValidateSession(t.Context(), tc.sessionID)

			require.ErrorIs(t, err, tc.expectErr)
		})
	}
}
