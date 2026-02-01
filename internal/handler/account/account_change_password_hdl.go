package account

import (
	"auth-service/gen/openapi"
	"auth-service/internal/apperrs"
	"auth-service/internal/biz/usecase/account"
	"auth-service/internal/handler/middlewares"
	"auth-service/pkg/transport/request"
	"auth-service/pkg/transport/response"

	"github.com/gin-gonic/gin"
)

func (hdl *accountHandler) ChangePassword(ctx *gin.Context) {
	authObj, err := middlewares.GetAuthObject(ctx)
	if err != nil {
		_ = ctx.Error(apperrs.Internal(err))
		return
	}

	payload, err := request.ParseAndValidateJSON[openapi.ChangePasswordJSONRequestBody](ctx)
	if err != nil {
		_ = ctx.Error(apperrs.InvalidInput(err))
		return
	}

	input := account.ChangePasswordInput{
		AccountID:   authObj.UserID,
		OldPassword: payload.OldPassword,
		NewPassword: payload.NewPassword,
	}
	if err := hdl.accountUC.ChangePassword(ctx.Request.Context(), input); err != nil {
		_ = ctx.Error(err)
		return
	}

	response.NoContent(ctx)
}
