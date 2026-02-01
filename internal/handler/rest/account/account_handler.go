package account

import (
	"auth-service/internal/usecase/port"

	"github.com/gin-gonic/gin"
)

type AccountHandler interface {
	CreateAccount(ctx *gin.Context)
	ChangePassword(ctx *gin.Context)
}

type accountHandler struct {
	accountUC port.AccountUsecase
	sessionUC port.SessionUsecase
}

func NewAccountHandler(accountUC port.AccountUsecase, sessionUC port.SessionUsecase) AccountHandler {
	return &accountHandler{
		accountUC: accountUC,
		sessionUC: sessionUC,
	}
}
