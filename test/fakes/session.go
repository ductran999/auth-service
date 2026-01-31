package fakes

import (
	"auth-service/internal/domain/sessionmodel"

	"github.com/google/uuid"
)

func Session() *sessionmodel.Session {
	return &sessionmodel.Session{
		ID:        uuid.MustParse("9f3c6d2a-4b7f-4c8a-9f6a-2e3d8b6f4a91"),
		AccountID: FakeAccount().ID,
	}
}
