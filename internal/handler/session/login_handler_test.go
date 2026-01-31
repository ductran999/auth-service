package session_test

import (
	"auth-service/gen/openapi"
	"auth-service/internal/handler/middlewares"
	"auth-service/internal/handler/session"
	mockbuilder "auth-service/test/mock-builder"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DucTran999/shared-pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func NewAuthSessionHandlerUT(t *testing.T, builder *mockbuilder.UsecaseBuilderContainer) session.SessionAuthHandler {
	t.Helper()
	return session.NewSessionAuthHandler(builder.AuthSessionUC.GetInstance())
}

func TestLoginAccount(t *testing.T) {
	tests := []struct {
		name           string
		setupUT        func(t *testing.T) session.SessionAuthHandler
		setupPayload   func(t *testing.T) []byte
		tokenKey       string
		expectedStatus int
	}{
		{
			name: "invalid json payload",
			setupUT: func(t *testing.T) session.SessionAuthHandler {
				builder := mockbuilder.NewUsecaseBuilderContainer(t)
				return NewAuthSessionHandlerUT(t, builder)
			},
			setupPayload: func(t *testing.T) []byte {
				// Invalid JSON
				return []byte(`{invalid`)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid credentials",
			setupUT: func(t *testing.T) session.SessionAuthHandler {
				builder := mockbuilder.NewUsecaseBuilderContainer(t)
				builder.AuthSessionUC.LoginInvalidCredentials()
				return NewAuthSessionHandlerUT(t, builder)
			},
			tokenKey: "session_id",
			setupPayload: func(t *testing.T) []byte {
				req := openapi.LoginAccountJSONRequestBody{
					Email:    "wrong@example.com",
					Password: "wrongpass",
				}
				b, _ := json.Marshal(req)
				return b
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "login success",
			setupUT: func(t *testing.T) session.SessionAuthHandler {
				builder := mockbuilder.NewUsecaseBuilderContainer(t)
				builder.AuthSessionUC.LoginSuccess()
				return NewAuthSessionHandlerUT(t, builder)
			},
			tokenKey: "session",
			setupPayload: func(t *testing.T) []byte {
				req := openapi.LoginAccountJSONRequestBody{
					Email:    "user@example.com",
					Password: "validPass123!",
				}
				b, _ := json.Marshal(req)
				return b
			},
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
			router.Use(middlewares.ErrorLogger(logger))
			router.POST("/api/v1/login", handler.LoginAccount)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(tc.setupPayload(t)))
			req.Header.Set("Content-Type", "application/json")

			// Optional: simulate existing session cookie
			req.AddCookie(&http.Cookie{
				Name:  tc.tokenKey,
				Value: "mock",
			})

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}
