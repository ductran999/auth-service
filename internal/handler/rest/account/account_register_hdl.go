package account

import (
	"auth-service/gen/openapi"
	"auth-service/internal/apperrs"
	"auth-service/internal/usecase/dto"
	"auth-service/pkg/transport/request"
	"auth-service/pkg/transport/response"

	"github.com/gin-gonic/gin"
)

// CreateAccount handles the HTTP request to register a new account.
func (hdl *accountHandler) CreateAccount(ctx *gin.Context) {
	payload, err := request.ParseAndValidateJSON[openapi.CreateAccountJSONRequestBody](ctx)
	if err != nil {
		_ = ctx.Error(apperrs.InvalidInput(err))
		return
	}

	input := dto.RegisterInput{
		Email:    payload.Email,
		Password: payload.Password,
	}

	account, err := hdl.accountUC.Register(ctx.Request.Context(), input)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	resp := openapi.Account{
		UserId:    account.ID,
		Email:     account.Email,
		Role:      account.Role,
		CreatedAt: &account.CreatedAt,
		UpdatedAt: &account.UpdatedAt,
	}

	response.OK(ctx, resp, "register new account successfully!")
}
