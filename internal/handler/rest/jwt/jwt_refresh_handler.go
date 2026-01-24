package jwt

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (hdl *jwtAuthHandler) RefreshToken(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie(RefreshTokenKey)
	if err != nil {
		hdl.UnauthorizeErrorResponse(ctx, APIVersion2, http.StatusText(http.StatusUnauthorized))
		return
	}

	tokens, err := hdl.authUC.RefreshToken(ctx, refreshToken)
	if err != nil {
		hdl.UnauthorizeErrorResponse(ctx, APIVersion2, err.Error())
		return
	}

	hdl.responseLoginJWTSuccess(ctx, tokens)
}
