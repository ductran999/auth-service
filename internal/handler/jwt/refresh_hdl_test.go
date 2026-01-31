package jwt_test

import (
	httpserver "auth-service/internal/app/server/http"
	"auth-service/internal/handler/jwt"
	"auth-service/internal/handler/middlewares"
	mockbuilder "auth-service/test/mock-builder"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DucTran999/shared-pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
)

func TestRefreshToken(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		mustSetCookie  bool
		setupUT        func(t *testing.T) jwt.JWTAuthHandler
		expectedStatus int
	}{
		{
			name: "missing refresh token in cookie",
			setupUT: func(t *testing.T) jwt.JWTAuthHandler {
				builder := mockbuilder.NewUsecaseBuilderContainer(t)
				return NewAuthJWTHandlerUT(t, builder)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "invalid refresh token usecase returns error",
			setupUT: func(t *testing.T) jwt.JWTAuthHandler {
				builder := mockbuilder.NewUsecaseBuilderContainer(t)
				builder.AuthJwtUC.RefreshTokenError()
				return NewAuthJWTHandlerUT(t, builder)
			},
			mustSetCookie:  true,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "refresh token success",
			setupUT: func(t *testing.T) jwt.JWTAuthHandler {
				builder := mockbuilder.NewUsecaseBuilderContainer(t)
				builder.AuthJwtUC.RefreshTokenSuccess()
				return NewAuthJWTHandlerUT(t, builder)
			},
			mustSetCookie:  true,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			handler := tc.setupUT(t)

			gin.SetMode(gin.TestMode)
			router := gin.New()
			logger, _ := logger.NewLogger(logger.Config{Environment: "staging"})
			err := httpserver.SetupValidator()
			require.NoError(t, err)
			router.Use(middlewares.ErrorLogger(logger))
			router.POST("/api/v2/refresh", handler.RefreshToken)

			req := httptest.NewRequest(http.MethodPost, "/api/v2/refresh", nil)
			req.Header.Set("Content-Type", "application/json")
			if tc.mustSetCookie {
				req.AddCookie(&http.Cookie{
					Name:  "refresh_token",
					Value: "mock-refresh-token",
					Path:  "/",
				})
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}
