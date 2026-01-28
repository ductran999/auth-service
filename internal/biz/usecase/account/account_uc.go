package account

import (
	"auth-service/internal/biz/ports/repositories"
	"auth-service/internal/domain/accountmodel"
	"auth-service/pkg/hasher"
	"context"
)

// AccountUseCase defines the business logic for managing user accounts.
type AccountUsecase interface {
	// Register creates a new user account with the provided information.
	// It typically includes validation, password hashing, and persistence logic.
	Register(ctx context.Context, input RegisterInput) (*accountmodel.Account, error)

	// ChangePassword change password for user when old password are match
	ChangePassword(ctx context.Context, input ChangePasswordInput) error
}

type accountUsecase struct {
	*registerUsecase
	*changePasswordUsecase
}

func NewAccountUseCase(hasher hasher.Hasher, accountRepo repositories.AccountRepo) AccountUsecase {
	return &accountUsecase{
		registerUsecase: &registerUsecase{
			hasher:      hasher,
			accountRepo: accountRepo,
		},
		changePasswordUsecase: &changePasswordUsecase{
			hasher:      hasher,
			accountRepo: accountRepo,
		},
	}
}
