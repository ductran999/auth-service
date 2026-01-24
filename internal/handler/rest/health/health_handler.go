package health

import (
	"auth-service/gen/openapi"
	"auth-service/pkg/transport/response"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthHandler interface {
	CheckLiveness(ctx *gin.Context)
}

type healthHandler struct {
	serviceVersion string // version release
	startTime      time.Time
}

func NewHealthHandler(serviceVersion string) HealthHandler {
	return &healthHandler{
		serviceVersion: serviceVersion,
		startTime:      time.Now(),
	}
}

func (h *healthHandler) CheckLiveness(ctx *gin.Context) {
	uptime := int64(time.Since(h.startTime).Seconds())
	respBody := openapi.HealthResponse{
		Status:    openapi.HealthResponseStatusHealthy,
		Uptime:    &uptime,
		Timestamp: time.Now(),
		Version:   &h.serviceVersion,
	}

	response.OK(ctx, respBody)
}
