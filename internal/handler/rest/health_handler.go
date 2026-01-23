package rest

import (
	"auth-service/gen/openapi"
	"net/http"
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

	response := openapi.HealthResponse{
		Status:    openapi.HealthResponseStatusHealthy,
		Timestamp: time.Now().UTC(),
		Uptime:    &uptime,
		Version:   &h.serviceVersion,
	}

	ctx.JSON(http.StatusOK, response)
}
