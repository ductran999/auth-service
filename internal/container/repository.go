package container

import (
	"auth-service/internal/infra/storage/postgresql"
	"auth-service/internal/infra/storage/session"
	"auth-service/internal/usecase/auth/ports"
	"auth-service/internal/usecase/port"
)

type repositories struct {
	account port.AccountRepo
	session ports.SessionRepository
}

func (c *Container) initRepositories() {
	c.repositories = &repositories{
		account: postgresql.NewAccountRepo(c.AuthDB.DB()),
		session: session.NewSessionRepository(c.AuthDB.DB()),
	}
}
