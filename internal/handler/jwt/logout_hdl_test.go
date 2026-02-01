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

func TestLogoutJWT(t *testing.T) {
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
			name: "revoke token got internal error",
			setupUT: func(t *testing.T) jwt.JWTAuthHandler {
				builder := mockbuilder.NewUsecaseBuilderContainer(t)
				builder.AuthJwtUC.RevokeRefreshTokenErr()
				return NewAuthJWTHandlerUT(t, builder)
			},
			mustSetCookie:  true,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "logout success",
			setupUT: func(t *testing.T) jwt.JWTAuthHandler {
				builder := mockbuilder.NewUsecaseBuilderContainer(t)
				builder.AuthJwtUC.RevokeRefreshTokenSuccess()
				return NewAuthJWTHandlerUT(t, builder)
			},
			mustSetCookie:  true,
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			handler := tc.setupUT(t)

			// Setup handler with mock
			gin.SetMode(gin.TestMode)
			router := gin.New()
			logger, _ := logger.NewLogger(logger.Config{Environment: "staging"})
			err := httpserver.SetupValidator()
			require.NoError(t, err)
			router.Use(middlewares.ErrorLogger(logger))
			router.POST("/api/v2/logout", handler.LogoutJWT)

			// Make request
			req := httptest.NewRequest(http.MethodPost, "/api/v2/logout", nil)
			req.Header.Set("Content-Type", "application/json")
			if tc.mustSetCookie {
				req.AddCookie(&http.Cookie{
					Name:  "refresh_token",
					Value: "mock-refresh-token",
					Path:  "/",
				})
			}

			// Setup response recorder
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}
