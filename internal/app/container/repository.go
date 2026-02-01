package container

import (
	"auth-service/internal/biz/ports/repositories"
	"auth-service/internal/infra/account"
	"auth-service/internal/infra/session"
)

type repos struct {
	account repositories.AccountRepo
	session repositories.SessionRepository
}

func (c *Container) initRepositories() {
	c.repos = &repos{
		account: account.NewAccountRepo(c.AuthDB.DB()),
		session: session.NewSessionRepository(c.AuthDB.DB()),
	}
}
