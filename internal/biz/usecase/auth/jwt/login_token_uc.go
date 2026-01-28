package jwt

import (
	"auth-service/internal/biz/ports/auth"
	"auth-service/internal/biz/usecase/auth/credential"
	"auth-service/internal/domain/authmodel"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	AccessTokenLifetime  = 15 * time.Minute
	RefreshTokenLifetime = 7 * 24 * time.Hour
)

type LoginJWTInput struct {
	Email     string
	Password  string
	IP        string
	UserAgent string
}

type loginWithJWTUsecase struct {
	tokenService auth.TokenService
	tokenStore   auth.TokenStore

	credVerifier *credential.CredentialVerifier
}

// Login authenticates the user and returns a pair of JWT tokens (access + refresh).
func (uc *loginWithJWTUsecase) Login(ctx context.Context, input LoginJWTInput) (*authmodel.TokenPairs, error) {
	// Verify user credentials
	account, err := uc.credVerifier.Verify(ctx, input.Email, input.Password)
	if err != nil {
		return nil, fmt.Errorf("auth: failed to verify credentials: %w", err)
	}

	// Prepare token metadata
	jti := uuid.NewString()
	issuedAt := time.Now()

	// Generate access and refresh tokens
	tokenPairs, err := uc.tokenService.SignPairs(jti, issuedAt, account)
	if err != nil {
		return nil, fmt.Errorf("auth: failed to sign token pairs: %w", err)
	}

	// Store Device Session
	deviceSession := authmodel.DeviceSession{
		JTI:       jti,
		AccountID: account.ID.String(),
		UserAgent: input.UserAgent,
		IP:        input.IP,
		SignAt:    issuedAt,
		ExpiresAt: issuedAt.Add(RefreshTokenLifetime),
	}
	if err := uc.tokenStore.Save(ctx, deviceSession); err != nil {
		return nil, fmt.Errorf("auth: failed to cache refresh token: %w", err)
	}

	return tokenPairs, nil
}
