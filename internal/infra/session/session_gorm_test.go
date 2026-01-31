package session_test

import (
	"errors"
	"testing"

	"auth-service/internal/domain/sessionmodel"
	"auth-service/internal/infra/session"
	mockbuilder "auth-service/test/mock-builder"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSessionRepo_Create_DBError(t *testing.T) {
	gormDB, mock, cleanup := mockbuilder.NewMockGormDB(t)
	defer cleanup()

	repo := session.NewSessionRepository(gormDB)

	session := &sessionmodel.Session{
		ID:        uuid.New(),
		AccountID: uuid.New(),
	}

	// Setup mock to simulate DB error
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "sessions"`).
		WithArgs(
			session.AccountID,
			sqlmock.AnyArg(), // ip_address
			sqlmock.AnyArg(), // user_agent
			sqlmock.AnyArg(), // last_seen_at
			sqlmock.AnyArg(), // revoked_at
			session.ID,       // id
		).
		WillReturnError(errors.New("simulated db error"))
	mock.ExpectRollback()

	// Act
	err := repo.Create(t.Context(), session)

	// Assert
	require.Error(t, err)
	require.Contains(t, err.Error(), "simulated db error")

	// Check mock expectations
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSessionRepo_Revoke_DBError(t *testing.T) {
	t.Parallel()
	db, mock, cleanup := mockbuilder.NewMockGormDB(t)
	defer cleanup()

	sessionID := "11111111-1111-1111-1111-111111111111"

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "sessions"`).
		WithArgs(sqlmock.AnyArg(), sessionID).
		WillReturnError(errors.New("simulated update failure"))
	mock.ExpectRollback()

	repo := session.NewSessionRepository(db)
	err := repo.Revoke(t.Context(), sessionID)

	require.Error(t, err)
	require.Contains(t, err.Error(), "simulated update failure")
	require.NoError(t, mock.ExpectationsWereMet())
}
