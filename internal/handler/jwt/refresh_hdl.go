package jwt

import (
	"auth-service/internal/apperrs"
	"auth-service/pkg/transport/response"

	"github.com/gin-gonic/gin"
)

func (hdl *jwtAuthHandler) RefreshToken(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie(RefreshTokenKey)
	if err != nil {
		_ = ctx.Error(apperrs.Unauthorized(err))
		return
	}

	tokens, err := hdl.authUC.RefreshToken(ctx, refreshToken)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	response.OK(ctx, tokens)
}
