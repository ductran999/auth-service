package container

import (
	"auth-service/internal/handler/rest"
	"auth-service/internal/handler/rest/health"
	jwtAuthHandler "auth-service/internal/handler/rest/jwt"
)

type RestHandler struct {
	rest.SessionAuthHandler
	jwtAuthHandler.JWTAuthHandler
	rest.AccountHandler
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
