package middlewares

import (
	"auth-service/internal/apperrs"
	"auth-service/pkg/transport/response"

	"github.com/DucTran999/shared-pkg/logger"
	"github.com/gin-gonic/gin"
)

func ErrorLogger(logger logger.ILogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		appErr := apperrs.ToAppError(err)

		status := apperrs.HTTPStatus(appErr)

		logger.Error(appErr.Cause.Error())

		response.ErrorDetail(c, status, appErr.Error())
	}
}
