package auth

import (
	"auth-service/internal/biz/ports/auth"
	"auth-service/internal/domain/accountmodel"
	"auth-service/internal/domain/authmodel"
	"context"
	"fmt"
	"time"

	"github.com/DucTran999/jwtkit"
	"github.com/golang-jwt/jwt/v5"
)

const (
	AccessTokenLifetime  = 15 * time.Minute
	RefreshTokenLifetime = 7 * 24 * time.Hour
)

type tokenClaimsParams struct {
	signAt     time.Time
	exp        time.Duration
	jti        string
	account    *accountmodel.Account
	includeJTI bool
}

type jwtSigner struct {
	signer jwtkit.JWT
}

func NewJWTSigner(signer jwtkit.JWT) auth.TokenService {
	return &jwtSigner{
		signer: signer,
	}
}

func (js *jwtSigner) Sign(claims authmodel.TokenClaims) (string, error) {
	return js.signer.Sign(claims)
}

func (js *jwtSigner) SignPairs(jti string, signAt time.Time, account *accountmodel.Account) (*authmodel.TokenPairs, error) {
	// Access token claims
	accessClaims := js.buildClaims(tokenClaimsParams{
		signAt:  signAt,
		exp:     AccessTokenLifetime,
		account: account,
	})

	accessToken, err := js.signer.Sign(accessClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Refresh token claims
	refreshClaims := js.buildClaims(tokenClaimsParams{
		signAt:     signAt,
		jti:        jti,
		exp:        RefreshTokenLifetime,
		account:    account,
		includeJTI: true,
	})

	refreshToken, err := js.signer.Sign(refreshClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return &authmodel.TokenPairs{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (js *jwtSigner) VerifyRefreshToken(context.Context, string) (*authmodel.TokenClaims, error) {
	return &authmodel.TokenClaims{}, nil
}

func (uc *jwtSigner) buildClaims(params tokenClaimsParams) authmodel.TokenClaims {
	claims := authmodel.TokenClaims{
		Email: params.account.Email,
		Role:  params.account.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   params.account.ID.String(),
			IssuedAt:  jwt.NewNumericDate(params.signAt),
			ExpiresAt: jwt.NewNumericDate(params.signAt.Add(params.exp)),
		},
	}

	if params.includeJTI {
		claims.ID = params.jti
	}

	return claims
}
