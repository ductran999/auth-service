package jwt

import (
	"auth-service/gen/openapi"
	"auth-service/internal/apperrs"
	"auth-service/internal/domain/authmodel"
	"auth-service/pkg/transport/response"
	"net/http"
	"time"

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

func (hdl *jwtAuthHandler) responseLoginJWTSuccess(ctx *gin.Context, tokens *authmodel.TokenPairs) {
	// Determine environment is secure or not
	secure := ctx.Request.Header.Get("X-Forwarded-Proto") == "https" || ctx.Request.TLS != nil

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     RefreshTokenKey,
		Value:    tokens.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(time.Second),
	})

	resp := openapi.AccessToken{
		AccessToken: tokens.AccessToken,
	}

	response.OK(ctx, resp)
}
