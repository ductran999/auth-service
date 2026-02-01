package jwt

import (
	"auth-service/gen/openapi"
	"auth-service/internal/apperrs"
	"auth-service/internal/biz/usecase/auth/jwt"
	"auth-service/internal/domain/authmodel"
	"auth-service/pkg/transport/request"
	"auth-service/pkg/transport/response"
	"net/http"
	"time"

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
		Expires:  time.Now().Add(time.Hour),
	})

	resp := openapi.AccessToken{
		AccessToken: tokens.AccessToken,
	}

	response.OK(ctx, resp)
}
