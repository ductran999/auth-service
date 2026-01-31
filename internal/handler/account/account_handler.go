package account

import (
	"auth-service/internal/biz/usecase/account"

	"github.com/gin-gonic/gin"
)

type AccountHandler interface {
	CreateAccount(ctx *gin.Context)
	ChangePassword(ctx *gin.Context)
}

type accountHandler struct {
	accountUC account.AccountUsecase
}

func NewAccountHandler(accountUC account.AccountUsecase) AccountHandler {
	return &accountHandler{
		accountUC: accountUC,
	}
}
