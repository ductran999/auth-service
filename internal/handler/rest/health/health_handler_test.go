package health_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"auth-service/internal/handler/rest/health"
	"auth-service/internal/handler/rest/middlewares"

	"github.com/DucTran999/shared-pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckLiveness(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	logger, err := logger.NewLogger(logger.Config{
		Environment: "unittest",
	})
	require.NoError(t, err)
	router := gin.New()
	router.Use(middlewares.ErrorLogger(logger))

	const version = "v1.2.3"
	handler := health.NewHealthHandler(version)

	router.GET("/livez", handler.CheckLiveness)

	// Make request
	req := httptest.NewRequest(http.MethodGet, "/livez", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
}
