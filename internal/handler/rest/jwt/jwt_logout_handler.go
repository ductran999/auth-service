package jwt

import (
	"auth-service/gen/openapi"
	"auth-service/internal/usecase/dto"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (hdl *jwtAuthHandler) LogoutJWT(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie(RefreshTokenKey)
	if err != nil {
		hdl.UnauthorizeErrorResponse(ctx, APIVersion2, http.StatusText(http.StatusUnauthorized))
		return
	}

	if err := hdl.authUC.RevokeRefreshToken(ctx, refreshToken); err != nil {
		hdl.UnauthorizeErrorResponse(ctx, APIVersion2, http.StatusText(http.StatusUnauthorized))
		return
	}

	// Always clear the cookie
	ctx.SetCookie(RefreshTokenKey, "", -1, "/", "", true, true)

	// Always respond with 204 No Content
	hdl.NoContentResponse(ctx)
}

func (hdl *jwtAuthHandler) responseLoginJWTSuccess(ctx *gin.Context, tokens *dto.TokenPairs) {
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

	resp := openapi.LoginJWTResponse{
		Success: true,
		Version: APIVersion2,
		Data: openapi.AccessToken{
			AccessToken: tokens.AccessToken,
		},
	}

	ctx.JSON(http.StatusOK, resp)
}
