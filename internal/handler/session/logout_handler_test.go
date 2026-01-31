package session_test

import (
	"auth-service/internal/handler/session"
	mockbuilder "auth-service/test/mock-builder"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLogoutAccount(t *testing.T) {
	tests := []struct {
		name           string
		mustSetCookie  bool
		setupUT        func(t *testing.T) session.SessionAuthHandler
		expectedStatus int
	}{
		{
			name:          "logout success with cookie",
			mustSetCookie: true,
			setupUT: func(t *testing.T) session.SessionAuthHandler {
				builder := mockbuilder.NewUsecaseBuilderContainer(t)
				builder.AuthSessionUC.LogoutSuccess()
				return NewAuthSessionHandlerUT(t, builder)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:          "logout error but still 204",
			mustSetCookie: true,
			setupUT: func(t *testing.T) session.SessionAuthHandler {
				builder := mockbuilder.NewUsecaseBuilderContainer(t)
				builder.AuthSessionUC.LogoutError()
				return NewAuthSessionHandlerUT(t, builder)
			},
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			handler := tc.setupUT(t)

			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.POST("/api/v1/logout", handler.LogoutAccount)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/logout", nil)
			if tc.mustSetCookie {
				req.AddCookie(&http.Cookie{
					Name:  "session_id",
					Value: "mock-session-id",
					Path:  "/",
				})
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)

			// Cookie must be cleared in response
			cookies := w.Result().Cookies()
			foundCleared := false
			for _, c := range cookies {
				if c.Name == "session_id" && c.MaxAge < 0 {
					foundCleared = true
				}
			}

			assert.True(t, foundCleared, "session_id cookie should be cleared")
		})
	}
}
