package middlewares

import (
	"time"

	"github.com/DucTran999/shared-pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RequestLogger(log logger.ILogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		log.Info("http request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", status),
			zap.String("latency", latency.String()),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}
