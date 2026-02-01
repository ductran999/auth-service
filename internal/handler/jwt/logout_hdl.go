package jwt

import (
	"auth-service/internal/apperrs"
	"auth-service/pkg/transport/response"

	"github.com/gin-gonic/gin"
)

func (hdl *jwtAuthHandler) LogoutJWT(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie(RefreshTokenKey)
	if err != nil {
		_ = ctx.Error(apperrs.Unauthorized(err))
		return
	}

	if err := hdl.authUC.RevokeRefreshToken(ctx, refreshToken); err != nil {
		_ = ctx.Error(err)
		return
	}

	// Always clear the cookie
	ctx.SetCookie(RefreshTokenKey, "", -1, "/", "", true, true)

	// Always respond with 204 No Content
	response.NoContent(ctx)
}
