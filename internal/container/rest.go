package container

import (
	"auth-service/internal/handler/rest/account"
	health "auth-service/internal/handler/rest/health"
	"auth-service/internal/handler/rest/jwt"
	"auth-service/internal/handler/rest/session"
)

type RestHandler struct {
	session.SessionAuthHandler
	jwt.JWTAuthHandler
	account.AccountHandler
	health.HealthHandler
}

func (c *Container) initRestHandler() {
	c.RestHandler = &RestHandler{
		JWTAuthHandler:     c.handlers.jwtAuth,
		SessionAuthHandler: c.handlers.sessionAuth,
		AccountHandler:     c.handlers.account,
		HealthHandler:      c.handlers.health,
	}
}
