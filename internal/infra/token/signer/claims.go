package signer

import (
	"auth-service/internal/model"
	"time"
)

type tokenClaimsParams struct {
	signAt     time.Time
	exp        time.Duration
	jti        string
	account    *model.Account
	includeJTI bool
}
