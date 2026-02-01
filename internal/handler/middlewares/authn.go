package middlewares

import (
	"auth-service/internal/apperrs"
	"auth-service/internal/biz/usecase/auth/session"
	"auth-service/internal/domain/authmodel"
	"errors"
	"slices"

	"github.com/gin-gonic/gin"
)

type CtxKey string

const AuthCtxKey = "auth_obj"

var (
	mustAuth = []string{
		"/api/v1/account/password",
		"/api/v1/logout",
	}
)

func Authenticate(authUC session.AuthSessionUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if !slices.Contains(mustAuth, path) {
			c.Next()
			return
		}

		sessionID, err := c.Cookie("session_id")
		if err != nil {
			_ = c.Error(apperrs.Unauthorized(session.ErrInvalidSession))
			c.Abort()
			return
		}

		authObj, err := authUC.ValidateSession(c.Request.Context(), sessionID)
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}

		c.Set(AuthCtxKey, authObj)
		c.Next()
	}
}

func GetAuthObject(c *gin.Context) (*authmodel.AuthObj, error) {
	result, ok := c.Get(AuthCtxKey)
	if !ok {
		return nil, errors.New("missing auth object")
	}

	authObj, ok := result.(*authmodel.AuthObj)
	if !ok {
		return nil, errors.New("invalid auth object")
	}

	return authObj, nil
}
