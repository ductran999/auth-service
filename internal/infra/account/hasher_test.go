package account_test

import (
	account "auth-service/internal/infra/account"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHasher_HashAndCompare_OK(t *testing.T) {
	h := account.NewHasher()

	hash, err := h.Hash("secret")
	require.NoError(t, err)

	ok, err := h.ComparePasswordAndHash("secret", hash)
	require.NoError(t, err)
	require.True(t, ok)
}

func TestHasher_Compare_WrongPassword(t *testing.T) {
	h := account.NewHasher()

	hash, _ := h.Hash("secret")
	ok, err := h.ComparePasswordAndHash("wrong", hash)

	require.NoError(t, err)
	require.False(t, ok)
}
