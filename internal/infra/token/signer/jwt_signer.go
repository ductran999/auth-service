package signer

import (
	"auth-service/internal/model"
	authPort "auth-service/internal/usecase/auth/ports"
	"auth-service/internal/usecase/dto"
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

type jwtSigner struct {
	signer jwtkit.JWT
}

func NewJWTSigner(signer jwtkit.JWT) authPort.TokenService {
	return &jwtSigner{
		signer: signer,
	}
}

func (js *jwtSigner) Sign(claims model.TokenClaims) (string, error) {
	return js.signer.Sign(claims)
}

func (js *jwtSigner) SignPairs(jti string, signAt time.Time, account *model.Account) (*dto.TokenPairs, error) {
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

	return &dto.TokenPairs{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (js *jwtSigner) VerifyAccessToken(context.Context, string) (model.TokenClaims, error) {
	return model.TokenClaims{}, nil
}

func (js *jwtSigner) VerifyRefreshToken(context.Context, string) (model.TokenClaims, error) {
	return model.TokenClaims{}, nil
}

func (uc *jwtSigner) buildClaims(params tokenClaimsParams) model.TokenClaims {
	claims := model.TokenClaims{
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
