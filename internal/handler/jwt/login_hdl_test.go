package jwt_test

import (
	httpserver "auth-service/internal/app/server/http"
	"auth-service/internal/handler/jwt"
	"auth-service/internal/handler/middlewares"
	mockbuilder "auth-service/test/mock-builder"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DucTran999/shared-pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
)

func NewAuthJWTHandlerUT(t *testing.T, builder *mockbuilder.UsecaseBuilderContainer) jwt.JWTAuthHandler {
	return jwt.NewJWTAuthHandler(
		builder.AuthJwtUC.GetInstance(),
	)
}

func TestLoginWithJWT(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		setupUT        func(t *testing.T) jwt.JWTAuthHandler
		setupPayload   func(t *testing.T) []byte
		expectedStatus int
	}{
		{
			name: "missing email",
			setupUT: func(t *testing.T) jwt.JWTAuthHandler {
				builder := mockbuilder.NewUsecaseBuilderContainer(t)
				return NewAuthJWTHandlerUT(t, builder)
			},
			setupPayload: func(t *testing.T) []byte {
				payload := map[string]any{
					"password": "p@ssG0rk1234!",
				}
				jsonPayload, err := json.Marshal(payload)
				require.NoError(t, err)
				return jsonPayload
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "wrong credentials",
			setupUT: func(t *testing.T) jwt.JWTAuthHandler {
				builder := mockbuilder.NewUsecaseBuilderContainer(t)
				builder.AuthJwtUC.LoginErrWrongCredentials()
				return NewAuthJWTHandlerUT(t, builder)
			},
			setupPayload: func(t *testing.T) []byte {
				payload := map[string]any{
					"email":    "danial@example.com",
					"password": "p@ssG0rk1234!",
				}
				jsonPayload, err := json.Marshal(payload)
				require.NoError(t, err)
				return jsonPayload
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "login DB error",
			setupUT: func(t *testing.T) jwt.JWTAuthHandler {
				builder := mockbuilder.NewUsecaseBuilderContainer(t)
				builder.AuthJwtUC.LoginErrDB()
				return NewAuthJWTHandlerUT(t, builder)
			},
			setupPayload: func(t *testing.T) []byte {
				payload := map[string]any{
					"email":    "danial@example.com",
					"password": "p@ssG0rk1234!",
				}
				jsonPayload, err := json.Marshal(payload)
				require.NoError(t, err)
				return jsonPayload
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "login success",
			setupUT: func(t *testing.T) jwt.JWTAuthHandler {
				builder := mockbuilder.NewUsecaseBuilderContainer(t)
				builder.AuthJwtUC.LoginSuccess()
				return NewAuthJWTHandlerUT(t, builder)
			},
			setupPayload: func(t *testing.T) []byte {
				payload := map[string]any{
					"email":    "danial@example.com",
					"password": "p@ssG0rk1234!",
				}
				jsonPayload, err := json.Marshal(payload)
				require.NoError(t, err)
				return jsonPayload
			},
			expectedStatus: http.StatusOK,
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
			router.POST("/api/v2/login", handler.LoginWithJWT)

			// Make request
			req := httptest.NewRequest(
				http.MethodPost,
				"/api/v2/login",
				bytes.NewBuffer(tc.setupPayload(t)),
			)
			req.Header.Set("Content-Type", "application/json")

			// Setup response recorder
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}
