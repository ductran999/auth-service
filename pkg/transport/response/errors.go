package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	InternalErrorCode      = "INTERNAL_ERROR"
	BadRequestErrorCode    = "BAD_REQUEST"
	UnauthorizedErrorCode  = "UNAUTHORIZED"
	ForbiddenErrorCode     = "FORBIDDEN"
	NotFoundErrorCode      = "NOT_FOUND"
	MethodNotAllowedCode   = "METHOD_NOT_ALLOWED"
	ConflictErrorCode      = "CONFLICT"
	TooManyRequestsCode    = "TOO_MANY_REQUESTS"
	ServiceUnavailableCode = "SERVICE_UNAVAILABLE"
)

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errorResp struct {
	Success bool        `json:"success"`
	Error   errorDetail `json:"error"`
}

func ErrorDetail(ctx *gin.Context, httpStatus int, errMsg string) {
	respBody := errorResp{
		Success: false,
		Error: errorDetail{
			Code:    errorCodeFromStatus(httpStatus),
			Message: errMsg,
		},
	}

	ctx.JSON(httpStatus, respBody)
}

func errorCodeFromStatus(status int) string {
	switch status {
	case http.StatusBadRequest:
		return BadRequestErrorCode
	case http.StatusUnauthorized:
		return UnauthorizedErrorCode
	case http.StatusForbidden:
		return ForbiddenErrorCode
	case http.StatusNotFound:
		return NotFoundErrorCode
	case http.StatusMethodNotAllowed:
		return MethodNotAllowedCode
	case http.StatusConflict:
		return ConflictErrorCode
	case http.StatusTooManyRequests:
		return TooManyRequestsCode
	case http.StatusServiceUnavailable:
		return ServiceUnavailableCode
	default:
		return InternalErrorCode
	}
}
