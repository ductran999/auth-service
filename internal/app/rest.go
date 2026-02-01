package app

import (
	"errors"
	"fmt"
	"net/http"

	"auth-service/internal/app/container"
	httpServer "auth-service/internal/app/server/http"

	"github.com/DucTran999/shared-pkg/server"
)

func startHTTPServer(c *container.Container) (server.HttpServer, error) {
	srv, err := httpServer.NewHTTPServer(c, c.RestHandler)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize HTTP server: %w", err)
	}

	go func() {
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			c.Logger.Fatalf("[FATAL] HTTP server crashed: %v", err)
		}
	}()

	return srv, nil
}
