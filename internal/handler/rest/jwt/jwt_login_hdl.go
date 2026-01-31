package jwt

import (
	"auth-service/gen/openapi"
	"auth-service/internal/apperrs"
	"auth-service/internal/biz/usecase/auth/jwt"
	"auth-service/pkg/transport/request"

	"github.com/gin-gonic/gin"
)

func (hdl *jwtAuthHandler) LoginWithJWT(ctx *gin.Context) {
	payload, err := request.ParseAndValidateJSON[openapi.LoginAccountJSONRequestBody](ctx)
	if err != nil {
		_ = ctx.Error(apperrs.InvalidInput(err))
		return
	}

	input := jwt.LoginJWTInput{
		Email:     payload.Email,
		Password:  payload.Password,
		IP:        ctx.ClientIP(),
		UserAgent: ctx.Request.UserAgent(),
	}

	tokens, err := hdl.authUC.Login(ctx, input)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	hdl.responseLoginJWTSuccess(ctx, tokens)
}
