//go:build integration

package account_test

import (
	"testing"

	"auth-service/config"
	"auth-service/internal/app/container"
	"auth-service/internal/biz/usecase/account"
	"auth-service/internal/domain/accountmodel"
	accountInfra "auth-service/internal/infra/account"

	"github.com/DucTran999/dbkit"
	"github.com/stretchr/testify/require"
)

func SetupTestDB(t *testing.T) dbkit.Connection {
	cfg, err := config.LoadConfig(".test.env")
	require.NoError(t, err)

	c, err := container.NewContainer(cfg)
	require.NoError(t, err)

	db := c.AuthDB.DB()

	err = db.Exec("TRUNCATE sessions, accounts CASCADE").Error
	require.NoError(t, err)

	return c.AuthDB
}

func TestAccountRepo(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	ctx := t.Context()
	repo := accountInfra.NewAccountRepo(db.DB())
	hasher := accountInfra.NewHasher()

	hashedPassword, err := hasher.Hash("sTr0ngP@ssg0rk")
	require.NoError(t, err)

	acc := accountmodel.Account{
		Email:        "daniel@example.go",
		PasswordHash: hashedPassword,
	}

	err = repo.Create(ctx, &acc)
	require.NoError(t, err)
	require.NotEmpty(t, acc.ID)

	t.Run("find by ID found", func(t *testing.T) {
		found, err := repo.FindByID(ctx, acc.ID.String())
		require.NoError(t, err)
		require.NotNil(t, found)
		require.Equal(t, acc.Email, found.Email)
	})

	t.Run("find by email found", func(t *testing.T) {
		found, err := repo.FindByEmail(ctx, acc.Email)
		require.NoError(t, err)
		require.NotNil(t, found)
		require.Equal(t, acc.ID, found.ID)
	})

	t.Run("find by email not found", func(t *testing.T) {
		found, err := repo.FindByEmail(ctx, "notfound@example.com")
		require.ErrorIs(t, err, account.ErrAccountNotFound)
		require.Nil(t, found)
	})

	t.Run("find by ID not found", func(t *testing.T) {
		found, err := repo.FindByID(ctx, "8f5c6b1e-dc99-4e33-a8c0-3e58fba86a65")
		require.ErrorIs(t, err, account.ErrAccountNotFound)
		require.Nil(t, found)
	})

	t.Run("update password ", func(t *testing.T) {
		newPassword := "n3wP@ssW0rd"
		hashedNewPassword, err := hasher.Hash(newPassword)
		require.NoError(t, err)

		err = repo.UpdatePasswordHash(ctx, acc.ID.String(), hashedNewPassword)
		require.NoError(t, err)

		updated, err := repo.FindByID(ctx, acc.ID.String())
		require.NoError(t, err)

		ok, err := hasher.ComparePasswordAndHash(newPassword, updated.PasswordHash)
		require.NoError(t, err)
		require.True(t, ok)
	})
}
