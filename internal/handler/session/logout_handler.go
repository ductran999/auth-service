package session

import (
	"auth-service/pkg/transport/response"

	"github.com/gin-gonic/gin"
)

func (hdl *sessionAuthHandler) LogoutAccount(ctx *gin.Context) {
	sessionID, _ := ctx.Cookie("session_Id")
	_ = hdl.authUC.Logout(ctx, sessionID)

	// Always clear the cookie
	ctx.SetCookie(SessionKey, "", -1, "/", "", true, true)

	// Always respond with 204 No Content
	response.NoContent(ctx)
}
