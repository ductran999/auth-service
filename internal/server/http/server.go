package http

import (
	"auth-service/config"
	"auth-service/gen/openapi"

	"github.com/DucTran999/shared-pkg/server"
)

// NewHTTPServer creates a new HTTP server with injected dependencies.
func NewHTTPServer(cfg *config.EnvConfiguration, apiHandler openapi.ServerInterface) (server.HttpServer, error) {
	serverConf := server.ServerConfig{
		Host: cfg.Host,
		Port: cfg.Port,
	}

	err := SetupValidator()
	if err != nil {
		return nil, err
	}

	router, err := NewRouter(cfg.ServiceEnv, apiHandler)
	if err != nil {
		return nil, err
	}

	httpServer, err := server.NewGinHttpServer(router, serverConf)
	if err != nil {
		return nil, err
	}

	return httpServer, nil
}
