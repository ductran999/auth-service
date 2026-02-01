package http

import (
	"auth-service/gen/openapi"
	"auth-service/internal/handler/rest/middlewares"
	"fmt"

	"github.com/DucTran999/shared-pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type RunningEnvironment int

const (
	ProductionEnv RunningEnvironment = iota
	DevelopmentEnv
)

func (r RunningEnvironment) String() string {
	switch r {
	case DevelopmentEnv:
		return "dev"
	case ProductionEnv:
		return "prod"
	// Set to default value dev if env invalid
	default:
		return "dev"
	}
}

func SetupValidator() error {
	// binding custom validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := RegisterStrongPasswordValidators(v); err != nil {
			return fmt.Errorf("failed to register strong validator: %w", err)
		}
	}

	return nil
}

func NewRouter(serviceEnv string, logger logger.ILogger, h openapi.ServerInterface) (*gin.Engine, error) {
	if serviceEnv == ProductionEnv.String() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(
		gin.Logger(),
		gin.Recovery(),
		middlewares.ErrorLogger(logger),
	)

	openapi.RegisterHandlers(router, h)

	return router, nil
}
