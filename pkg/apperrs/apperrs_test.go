package apperrs_test

import (
	"auth-service/pkg/apperrs"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	ErrDatabaseConnectionFailed = errors.New("database connection failed")
)

func TestAppErrors_WrapError(t *testing.T) {
	err := apperrs.Internal(ErrDatabaseConnectionFailed)
	assert.ErrorIs(t, err, apperrs.ErrInternal)
}

func TestAppErrors_ToAppError(t *testing.T) {
	originalErr := apperrs.NotFound(errors.New("user not found"))
	appErr := apperrs.ToAppError(originalErr)

	require.ErrorIs(t, appErr, apperrs.ErrNotFound)
	assert.Equal(t, "user not found", appErr.Error())
}
