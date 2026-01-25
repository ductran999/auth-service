package response

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Success struct {
	Success bool   `json:"success"`
	Data    any    `json:"data"`
	Message string `json:"message,omitempty"`
}

func OK(ctx *gin.Context, data any, messages ...string) {
	ctx.JSON(http.StatusOK, Success{
		Success: true,
		Data:    data,
		Message: strings.Join(messages, ", "),
	})
}

func NoContent(ctx *gin.Context) {
	ctx.Status(http.StatusNoContent)
}
