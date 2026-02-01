package app

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
)

func Run() error {
	appCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Load application configuration
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("load config error: %w", err)
	}

	// Initialize container
	c, err := initContainer(cfg)
	if err != nil {
		return fmt.Errorf("init container error: %w", err)
	}
	defer c.Close()

	// start Rest server
	restSrv, err := startHTTPServer(c)
	if err != nil {
		return fmt.Errorf("start http server error: %w", err)
	}

	// gracefully shutdown
	waitForShutdown(appCtx, restSrv, c)

	return nil
}
