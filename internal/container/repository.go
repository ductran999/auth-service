package container

import (
	"auth-service/internal/infra/storage/postgresql"
	"auth-service/internal/usecase/port"
)

type repositories struct {
	account port.AccountRepo
	session port.SessionRepository
}

func (c *Container) initRepositories() {
	c.repositories = &repositories{
		account: postgresql.NewAccountRepo(c.AuthDB.DB()),
		session: postgresql.NewSessionRepository(c.AuthDB.DB()),
	}
}
