package fakes

import (
	"auth-service/internal/domain/accountmodel"
	"time"

	"github.com/google/uuid"
)

var fixedTime = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

func FakeAccount() *accountmodel.Account {
	return &accountmodel.Account{
		ID:           uuid.MustParse("9f3c6d2a-4b7f-4c8a-9f6a-2e3d8b6f4a91"),
		Email:        "unittest@example.com",
		IsActive:     true,
		IsVerified:   true,
		Role:         "Devops",
		CreatedAt:    fixedTime,
		UpdatedAt:    fixedTime,
		PasswordHash: "hashed-password",
	}
}

func FakeAccountInactive() *accountmodel.Account {
	return &accountmodel.Account{
		ID:           uuid.MustParse("9f3c6d2a-4b7f-4c8a-9f6a-2e3d8b6f4a91"),
		Email:        "unittest@example.com",
		IsActive:     false,
		IsVerified:   true,
		Role:         "Devops",
		CreatedAt:    fixedTime,
		UpdatedAt:    fixedTime,
		PasswordHash: "hashed-password",
	}
}
