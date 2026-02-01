package account

import (
	"auth-service/internal/biz/usecase/account"
	"auth-service/internal/biz/usecase/port"

	"github.com/gin-gonic/gin"
)

type AccountHandler interface {
	CreateAccount(ctx *gin.Context)
	ChangePassword(ctx *gin.Context)
}

type accountHandler struct {
	accountUC account.AccountUsecase
	sessionUC port.SessionUsecase
}

func NewAccountHandler(accountUC account.AccountUsecase, sessionUC port.SessionUsecase) AccountHandler {
	return &accountHandler{
		accountUC: accountUC,
		sessionUC: sessionUC,
	}
}
