package account_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"auth-service/gen/openapi"
	"auth-service/internal/handler/rest/account"
	"auth-service/internal/handler/rest/middlewares"
	httpServer "auth-service/internal/server/http"
	mockbuilder "auth-service/test/mock-builder"

	"github.com/DucTran999/shared-pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
)

func NewAccountHandlerUT(t *testing.T, builder *mockbuilder.UsecaseBuilderContainer) account.AccountHandler {
	return account.NewAccountHandler(
		builder.AccountUC.GetInstance(),
		builder.SessionUC.GetInstance(),
	)
}

func TestCreateAccount(t *testing.T) {
	t.Parallel()

	type testcase struct {
		name           string
		setupUT        func(t *testing.T) account.AccountHandler
		setupPayload   func(t *testing.T) []byte
		expectedStatus int
	}

	tests := []testcase{
		{
			name: "missing email",
			setupUT: func(t *testing.T) account.AccountHandler {
				t.Helper()
				b := mockbuilder.NewUsecaseBuilderContainer(t)
				return NewAccountHandlerUT(t, b)
			},
			setupPayload: func(t *testing.T) []byte {
				t.Helper()
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
			name: "invalid email",
			setupUT: func(t *testing.T) account.AccountHandler {
				t.Helper()
				b := mockbuilder.NewUsecaseBuilderContainer(t)
				return NewAccountHandlerUT(t, b)
			},
			setupPayload: func(t *testing.T) []byte {
				t.Helper()
				payload := map[string]any{
					"email":    "invalidEmail.com",
					"password": "p@ssG0rk1234!",
				}
				jsonPayload, err := json.Marshal(payload)
				require.NoError(t, err)
				return jsonPayload
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid payload",
			setupUT: func(t *testing.T) account.AccountHandler {
				t.Helper()
				b := mockbuilder.NewUsecaseBuilderContainer(t)
				return NewAccountHandlerUT(t, b)
			},
			setupPayload: func(t *testing.T) []byte {
				t.Helper()
				payload := map[string]any{
					"email":    1234,
					"password": "p@ssG0rk1234!",
				}
				jsonPayload, err := json.Marshal(payload)
				require.NoError(t, err)
				return jsonPayload
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing password",
			setupUT: func(t *testing.T) account.AccountHandler {
				t.Helper()
				b := mockbuilder.NewUsecaseBuilderContainer(t)
				return NewAccountHandlerUT(t, b)
			},
			setupPayload: func(t *testing.T) []byte {
				t.Helper()
				payload := map[string]any{
					"email": "daniel@example.com",
				}
				jsonPayload, err := json.Marshal(payload)
				require.NoError(t, err)
				return jsonPayload
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "weak password",
			setupUT: func(t *testing.T) account.AccountHandler {
				t.Helper()
				b := mockbuilder.NewUsecaseBuilderContainer(t)
				return NewAccountHandlerUT(t, b)
			},
			setupPayload: func(t *testing.T) []byte {
				t.Helper()
				payload := map[string]any{
					"email":    "daniel@example.com",
					"password": "weak",
				}
				jsonPayload, err := json.Marshal(payload)
				require.NoError(t, err)
				return jsonPayload
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "email already registered",
			setupUT: func(t *testing.T) account.AccountHandler {
				t.Helper()
				b := mockbuilder.NewUsecaseBuilderContainer(t)
				b.AccountUC.RegisterConflictEmail()
				return NewAccountHandlerUT(t, b)
			},
			setupPayload: func(t *testing.T) []byte {
				t.Helper()
				payload := openapi.CreateAccountRequest{
					Email:    "test@example.com",
					Password: "p@ssG0rk1234!",
				}
				jsonPayload, err := json.Marshal(payload)
				require.NoError(t, err)
				return jsonPayload
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name: "register failed",
			setupUT: func(t *testing.T) account.AccountHandler {
				t.Helper()
				b := mockbuilder.NewUsecaseBuilderContainer(t)
				b.AccountUC.RegisterError()
				return NewAccountHandlerUT(t, b)
			},
			setupPayload: func(t *testing.T) []byte {
				t.Helper()
				payload := openapi.CreateAccountRequest{
					Email:    "test@example.com",
					Password: "p@ssG0rk1234!",
				}
				jsonPayload, err := json.Marshal(payload)
				require.NoError(t, err)
				return jsonPayload
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "register success",
			setupUT: func(t *testing.T) account.AccountHandler {
				t.Helper()
				b := mockbuilder.NewUsecaseBuilderContainer(t)
				b.AccountUC.RegisterSuccess()
				return NewAccountHandlerUT(t, b)
			},
			setupPayload: func(t *testing.T) []byte {
				t.Helper()
				payload := openapi.CreateAccountRequest{
					Email:    "test@example.com",
					Password: "p@ssG0rk1234!",
				}
				jsonPayload, err := json.Marshal(payload)
				require.NoError(t, err)
				return jsonPayload
			},
			expectedStatus: http.StatusOK,
		},
	}

	err := httpServer.SetupValidator()
	require.NoError(t, err)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			handler := tc.setupUT(t)

			// Setup handler with mock
			gin.SetMode(gin.TestMode)
			logger, _ := logger.NewLogger(logger.Config{Environment: "staging"})
			router := gin.New()
			router.Use(middlewares.ErrorLogger(logger))
			router.POST("/api/v1/register", handler.CreateAccount)

			// Setup response recorder
			req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBuffer(tc.setupPayload(t)))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}
