package http

import (
	"auth-service/gen/openapi"
	"auth-service/internal/app/container"

	"github.com/DucTran999/shared-pkg/server"
)

// NewHTTPServer creates a new HTTP server with injected dependencies.
func NewHTTPServer(ctn *container.Container, apiHandler openapi.ServerInterface) (server.HttpServer, error) {
	serverConf := server.ServerConfig{
		Host: ctn.AppConfig.Host,
		Port: ctn.AppConfig.Port,
	}

	err := SetupValidator()
	if err != nil {
		return nil, err
	}

	router, err := NewRouter(ctn, apiHandler)
	if err != nil {
		return nil, err
	}

	httpServer, err := server.NewGinHttpServer(router, serverConf)
	if err != nil {
		return nil, err
	}

	return httpServer, nil
}
