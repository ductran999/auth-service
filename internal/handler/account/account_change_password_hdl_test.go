package account_test

import (
	"auth-service/gen/openapi"
	"auth-service/internal/domain/authmodel"
	"auth-service/internal/handler/account"
	"auth-service/internal/handler/middlewares"
	mockbuilder "auth-service/test/mock-builder"
	"bytes"
	"encoding/json"

	"net/http"
	"net/http/httptest"
	"testing"

	httpserver "auth-service/internal/app/server/http"

	"github.com/DucTran999/shared-pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
)

func withAuthObj(authObj any) gin.HandlerFunc {
	return func(c *gin.Context) {
		if authObj != nil {
			c.Set(middlewares.AuthCtxKey, authObj)
		}
		c.Next()
	}
}

func TestAccountHandler_ChangePassword(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	endpoint := "/api/v1/account/password"

	tests := []struct {
		name           string
		authObj        any
		setupUT        func(t *testing.T) account.AccountHandler
		body           func(t *testing.T) []byte
		expectedStatus int
	}{
		{
			name:    "missing auth object",
			authObj: nil,
			setupUT: func(t *testing.T) account.AccountHandler {
				t.Helper()
				b := mockbuilder.NewUsecaseBuilderContainer(t)
				return NewAccountHandlerUT(t, b)
			},
			body: func(t *testing.T) []byte {
				t.Helper()
				b, _ := json.Marshal(`{}`)
				return b
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:    "invalid json payload",
			authObj: &authmodel.AuthObj{UserID: "acc-1"},
			setupUT: func(t *testing.T) account.AccountHandler {
				t.Helper()
				b := mockbuilder.NewUsecaseBuilderContainer(t)
				return NewAccountHandlerUT(t, b)
			},
			body: func(t *testing.T) []byte {
				t.Helper()
				b, _ := json.Marshal(`{invalid}`)
				return b
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:    "change password got error from usecase",
			authObj: &authmodel.AuthObj{UserID: "acc-1"},
			setupUT: func(t *testing.T) account.AccountHandler {
				t.Helper()
				b := mockbuilder.NewUsecaseBuilderContainer(t)
				b.AccountUC.ChangePassErrGotErrorDB()
				return NewAccountHandlerUT(t, b)
			},
			body: func(t *testing.T) []byte {
				t.Helper()
				b, err := json.Marshal(openapi.ChangePasswordRequest{
					OldPassword: "0ldP@ssGrork1234",
					NewPassword: "newPassWork123@23",
				})
				require.NoError(t, err)
				return b
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:    "change password success",
			authObj: &authmodel.AuthObj{UserID: "acc-1"},
			setupUT: func(t *testing.T) account.AccountHandler {
				t.Helper()
				b := mockbuilder.NewUsecaseBuilderContainer(t)
				b.AccountUC.ChangePasswordSuccess()
				return NewAccountHandlerUT(t, b)
			},
			body: func(t *testing.T) []byte {
				t.Helper()
				b, err := json.Marshal(openapi.ChangePasswordRequest{
					OldPassword: "0ldP@ssGrork1234",
					NewPassword: "newPassWork123@23",
				})
				require.NoError(t, err)
				return b
			},
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			handler := tc.setupUT(t)

			// Setup gin server
			logger, _ := logger.NewLogger(logger.Config{Environment: "staging"})
			router := gin.New()
			err := httpserver.SetupValidator()
			require.NoError(t, err)
			router.Use(middlewares.ErrorLogger(logger), withAuthObj(tc.authObj))
			router.POST(endpoint, handler.ChangePassword)

			// Setup response recorder
			req := httptest.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(tc.body(t)))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}
