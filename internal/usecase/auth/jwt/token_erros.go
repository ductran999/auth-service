package jwt

import "errors"

var (
	ErrInvalidRefreshToken = errors.New("token refresh token invalid")
)
