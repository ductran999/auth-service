package session

import (
	"auth-service/gen/openapi"
	"auth-service/internal/apperrs"
	"auth-service/internal/biz/usecase/auth/session"
	"auth-service/internal/domain/sessionmodel"
	"auth-service/pkg/transport/request"
	"auth-service/pkg/transport/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	SessionKey = "session_id"
)

type SessionAuthHandler interface {
	LoginAccount(ctx *gin.Context)
	LogoutAccount(ctx *gin.Context)
}

type sessionAuthHandler struct {
	authUC session.AuthSessionUsecase
}

func NewSessionAuthHandler(authUC session.AuthSessionUsecase) SessionAuthHandler {
	return &sessionAuthHandler{
		authUC: authUC,
	}
}

func (hdl *sessionAuthHandler) LoginAccount(ctx *gin.Context) {
	// Parse request body
	payload, err := request.ParseAndValidateJSON[openapi.LoginAccountJSONRequestBody](ctx)
	if err != nil {
		_ = ctx.Error(apperrs.InvalidInput(err))
		return
	}

	// Set to empty when cookie not found
	currentSessionID, err := ctx.Cookie(SessionKey)
	if err != nil {
		currentSessionID = ""
	}

	// Convert request to model
	loginInput := session.LoginInput{
		CurrentSessionID: currentSessionID,
		Email:            payload.Email,
		Password:         payload.Password,
		IP:               ctx.ClientIP(),
		UserAgent:        ctx.Request.UserAgent(),
	}

	// Authenticate user and create session
	session, err := hdl.authUC.Login(ctx.Request.Context(), loginInput)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	hdl.responseLoginSuccess(ctx, session)
}

func (hdl *sessionAuthHandler) responseLoginSuccess(ctx *gin.Context, session *sessionmodel.Session) {
	// Determine environment is secure or not
	secure := ctx.Request.Header.Get("X-Forwarded-Proto") == "https" || ctx.Request.TLS != nil

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     SessionKey,
		Value:    session.ID.String(),
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
	})

	resp := openapi.Account{
		UserId: session.AccountID,
	}

	response.OK(ctx, resp, "login successfully!")
}
