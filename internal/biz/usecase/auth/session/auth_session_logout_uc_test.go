package session_test

import (
	"auth-service/internal/biz/usecase/auth/session"
	mockbuilder "auth-service/test/mock-builder"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogout(t *testing.T) {
	type testcase struct {
		name      string
		sessionID string
		setup     func(t *testing.T) session.AuthSessionUsecase
		expectErr error
	}

	testTable := []testcase{
		{
			name:      "failed to revoke session",
			sessionID: mockbuilder.FakeSessionID.String(),
			setup: func(t *testing.T) session.AuthSessionUsecase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.SessionRepoBuilder.RevokeFailed()
				return NewAuthUseCaseUT(t, builders)
			},
			expectErr: mockbuilder.ErrUpdateSessionExpires,
		},
		{
			name:      "logout success",
			sessionID: mockbuilder.FakeSessionID.String(),
			setup: func(t *testing.T) session.AuthSessionUsecase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.SessionRepoBuilder.RevokeSuccess()
				return NewAuthUseCaseUT(t, builders)
			},
			expectErr: nil,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			sut := tc.setup(t)
			ctx := t.Context()

			err := sut.Logout(ctx, tc.sessionID)

			require.ErrorIs(t, err, tc.expectErr)
		})
	}
}
