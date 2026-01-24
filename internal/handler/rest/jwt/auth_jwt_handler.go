package jwt

import (
	"auth-service/gen/openapi"
	jwtAuthUC "auth-service/internal/usecase/auth/jwt"
	"auth-service/internal/usecase/dto"
	"auth-service/internal/usecase/port"
	"errors"

	"github.com/DucTran999/shared-pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	RefreshTokenKey = "refresh_token"
	RefreshTokenLifetime
)

type JWTAuthHandler interface {
	LoginWithJWT(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
	LogoutJWT(ctx *gin.Context)
}

type jwtAuthHandler struct {
	BaseHandler

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

func (hdl *jwtAuthHandler) LoginWithJWT(ctx *gin.Context) {
	// Parse request body
	payload, err := ParseAndValidateJSON[openapi.LoginAccountJSONRequestBody](ctx)
	if err != nil {
		hdl.BadRequestResponse(ctx, APIVersion2, err.Error())
		return
	}

	// Prepare input
	input := dto.LoginJWTInput{
		Email:     payload.Email,
		Password:  payload.Password,
		IP:        ctx.ClientIP(),
		UserAgent: ctx.Request.UserAgent(),
	}

	// Authenticate
	tokens, err := hdl.authUC.Login(ctx, input)
	if err != nil {
		if errors.Is(err, jwtAuthUC.ErrInvalidRefreshToken) {
			hdl.UnauthorizeErrorResponse(ctx, APIVersion2, err.Error())
		} else {
			hdl.logger.Error("failed to login with jwt", zap.String("error", err.Error()))
			hdl.ServerInternalErrResponse(ctx, APIVersion2)
		}
		return
	}

	hdl.responseLoginJWTSuccess(ctx, tokens)
}
