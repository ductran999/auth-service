package fakes

import (
	"auth-service/internal/domain/authmodel"

	"github.com/golang-jwt/jwt/v5"
)

const (
	fakeAT = "fake-access-token"
	fakeRT = "fake-refresh-token"
)

func FakeTokenPairs() *authmodel.TokenPairs {
	return &authmodel.TokenPairs{
		AccessToken:  fakeAT,
		RefreshToken: fakeRT,
	}
}

func FakeDeviceSession() *authmodel.DeviceSession {
	return &authmodel.DeviceSession{
		AccountID: FakeAccount().ID.String(),
	}
}

func FakeTokenClaims() *authmodel.TokenClaims {
	return &authmodel.TokenClaims{
		Email:            FakeAccount().Email,
		RegisteredClaims: jwt.RegisteredClaims{},
	}
}
