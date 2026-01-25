package container

import (
	"auth-service/internal/handler/rest/account"
	"auth-service/internal/handler/rest/health"
	"auth-service/internal/handler/rest/jwt"
	"auth-service/internal/handler/rest/session"
)

type handlers struct {
	sessionAuth session.SessionAuthHandler
	jwtAuth     jwt.JWTAuthHandler
	account     account.AccountHandler
	health      health.HealthHandler
}

func (c *Container) initHandlers() {
	c.handlers = &handlers{
		health: health.NewHealthHandler(c.AppConfig.ServiceVersion),

		sessionAuth: session.NewSessionAuthHandler(c.Logger, c.useCases.sessionAuth),
		jwtAuth:     jwt.NewJWTAuthHandler(c.Logger, c.useCases.jwtAuth),

		account: account.NewAccountHandler(c.useCases.account, c.useCases.session),
	}
}
