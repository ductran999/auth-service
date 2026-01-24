package container

import (
	"auth-service/internal/handler/rest"
	health "auth-service/internal/handler/rest/health"
	jwtAuth "auth-service/internal/handler/rest/jwt"
)

type handlers struct {
	sessionAuth rest.SessionAuthHandler
	jwtAuth     jwtAuth.JWTAuthHandler
	account     rest.AccountHandler
	health      health.HealthHandler
}

func (c *Container) initHandlers() {
	c.handlers = &handlers{
		health: health.NewHealthHandler(c.AppConfig.ServiceVersion),

		sessionAuth: rest.NewSessionAuthHandler(c.Logger, c.useCases.sessionAuth),
		jwtAuth:     jwtAuth.NewJWTAuthHandler(c.Logger, c.useCases.jwtAuth),

		account: rest.NewAccountHandler(c.Logger, c.useCases.account, c.useCases.session),
	}
}
