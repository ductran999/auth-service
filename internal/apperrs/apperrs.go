package apperrs

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrInternal     = errors.New("internal error")
	ErrInvalidInput = errors.New("invalid input")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
	ErrUnavailable  = errors.New("service unavailable")
)

type AppError struct {
	Err   error // sentinel error (ErrInternal, ErrNotFound, ...)
	Msg   string
	Cause error
}

func (e *AppError) Error() string {
	if e.Msg != "" {
		return e.Msg
	}
	return e.Err.Error()
}

func (e *AppError) Unwrap() error {
	if e.Cause != nil {
		return e.Cause
	}
	return e.Err
}

func ToAppError(err error) *AppError {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	return Internal(err)
}

func HTTPStatus(err error) int {
	switch {
	case errors.Is(err, ErrInvalidInput):
		return http.StatusBadRequest
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrConflict):
		return http.StatusConflict
	case errors.Is(err, ErrUnavailable):
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

// ===== Wrappers (preserve cause) =====

func newAppError(errType error, msg string, cause error) *AppError {
	causeWrapped := fmt.Errorf("%w: %w", errType, cause)
	appErr := &AppError{
		Err:   errType,
		Msg:   msg,
		Cause: causeWrapped,
	}

	return appErr
}

func Internal(cause error) *AppError {
	return newAppError(ErrInternal, "internal error", cause)
}

func InvalidInput(cause error) *AppError {
	return newAppError(ErrInvalidInput, cause.Error(), cause)
}

func Unauthorized(err error) *AppError {
	return newAppError(ErrUnauthorized, err.Error(), err)
}

func Forbidden(err error) *AppError {
	return newAppError(ErrForbidden, err.Error(), err)
}

func NotFound(err error) *AppError {
	return newAppError(ErrNotFound, err.Error(), err)
}

func Conflict(err error) *AppError {
	return newAppError(ErrConflict, err.Error(), err)
}

func Unavailable(cause error) *AppError {
	return newAppError(ErrUnavailable, ErrUnavailable.Error(), cause)
}
