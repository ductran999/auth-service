package jwt

import (
	"auth-service/internal/biz/usecase/auth/jwt"

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
	authUC jwt.AuthJWTUsecase
}

// NewJWTAuthHandler creates a new JWT authentication handler.
func NewJWTAuthHandler(logger logger.ILogger, authUC jwt.AuthJWTUsecase) JWTAuthHandler {
	return &jwtAuthHandler{
		logger: logger,
		authUC: authUC,
	}
}
