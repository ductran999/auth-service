package jwt

import (
	"auth-service/internal/usecase/port"

	"github.com/DucTran999/shared-pkg/logger"
	"github.com/gin-gonic/gin"
)

const (
	RefreshTokenKey = "refresh_token"
)

type JWTAuthHandler interface {
	LoginWithJWT(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
	LogoutJWT(ctx *gin.Context)
}

type jwtAuthHandler struct {
	logger logger.ILogger
	authUC port.AuthJWTUsecase
}

// NewJWTAuthHandler creates a new JWT authentication handler.
func NewJWTAuthHandler(logger logger.ILogger, authUC port.AuthJWTUsecase) JWTAuthHandler {
	return &jwtAuthHandler{
		logger: logger,
		authUC: authUC,
	}
}
