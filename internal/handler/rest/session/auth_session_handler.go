package session

import (
	"net/http"

	"auth-service/gen/openapi"
	"auth-service/internal/apperrs"
	"auth-service/internal/model"
	"auth-service/internal/usecase/dto"
	"auth-service/internal/usecase/port"
	"auth-service/pkg/transport/request"
	"auth-service/pkg/transport/response"

	"github.com/DucTran999/shared-pkg/logger"
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
	logger logger.ILogger
	authUC port.AuthSessionUsecase
}

func NewSessionAuthHandler(logger logger.ILogger, authUC port.AuthSessionUsecase) SessionAuthHandler {
	return &sessionAuthHandler{
		logger: logger,
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
	loginInput := dto.LoginInput{
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

func (hdl *sessionAuthHandler) LogoutAccount(ctx *gin.Context) {
	// Try to get session ID from cookie
	sessionID, err := ctx.Cookie(SessionKey)
	if err == nil {
		// Best-effort logout
		if err := hdl.authUC.Logout(ctx.Request.Context(), sessionID); err != nil {
			hdl.logger.Warn(err.Error())
		}
	}

	// Always clear the cookie
	ctx.SetCookie(SessionKey, "", -1, "/", "", true, true)

	// Always respond with 204 No Content
	response.NoContent(ctx)
}

func (hdl *sessionAuthHandler) responseLoginSuccess(ctx *gin.Context, session *model.Session) {
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
		Email:  session.Account.Email,
		Role:   session.Account.Role,
	}

	response.OK(ctx, resp, "login successfully!")
}
