package authmodel

import "github.com/golang-jwt/jwt/v5"

type TokenPairs struct {
	AccessToken  string
	RefreshToken string
}

type TokenClaims struct {
	jwt.RegisteredClaims

	// Custom fields
	Email string `json:"email"`
	Role  string `json:"role"`
}
