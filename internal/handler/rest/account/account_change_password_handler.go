package account

import (
	"auth-service/gen/openapi"
	"auth-service/internal/model"
	"auth-service/internal/usecase/dto"
	sessionUC "auth-service/internal/usecase/session"
	"auth-service/pkg/transport/request"
	"auth-service/pkg/transport/response"

	"github.com/gin-gonic/gin"
)

func (hdl *accountHandler) ChangePassword(ctx *gin.Context) {
	// Validate session from cookie
	session, err := hdl.validateSessionFromCookie(ctx)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	// Parse & validate input JSON
	payload, err := request.ParseAndValidateJSON[openapi.ChangePasswordJSONRequestBody](ctx)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	input := dto.ChangePasswordInput{
		AccountID:   session.AccountID.String(),
		OldPassword: payload.OldPassword,
		NewPassword: payload.NewPassword,
	}
	if err := hdl.accountUC.ChangePassword(ctx.Request.Context(), input); err != nil {
		_ = ctx.Error(err)
		return
	}

	response.NoContent(ctx)
}

func (hdl *accountHandler) validateSessionFromCookie(ctx *gin.Context) (*model.Session, error) {
	sessionID, err := ctx.Cookie("session_id")
	if err != nil {
		return nil, sessionUC.ErrSessionNotFound
	}

	return hdl.sessionUC.Validate(ctx.Request.Context(), sessionID)
}
