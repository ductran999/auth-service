package account

import (
	"auth-service/internal/apperrs"
	"auth-service/internal/biz/ports/repositories"
	"auth-service/internal/domain/accountmodel"
	"auth-service/pkg/hasher"
	"context"
	"errors"
)

type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerUsecase struct {
	hasher      hasher.Hasher
	accountRepo repositories.AccountRepo
}

func (uc *registerUsecase) Register(ctx context.Context, input RegisterInput) (*accountmodel.Account, error) {
	taken, err := uc.isEmailTaken(ctx, input.Email)
	if err != nil {
		return nil, apperrs.Internal(err)
	}
	if taken {
		return nil, apperrs.Conflict(ErrEmailExisted)
	}

	// Hash the password
	hashedPassword, err := uc.hasher.HashPassword(input.Password)
	if err != nil {
		return nil, apperrs.Internal(err)
	}

	// Bind input to domain model
	account := accountmodel.Account{
		Email:        input.Email,
		PasswordHash: hashedPassword,
	}

	// Persist the account
	if err := uc.accountRepo.Create(ctx, &account); err != nil {
		return nil, apperrs.Internal(err)
	}

	return &account, nil
}

// isEmailTaken checks if the provided email already exists in the system.
// Returns ErrEmailExisted if a duplicate is found, or a repository error if any occurs.
func (uc *registerUsecase) isEmailTaken(ctx context.Context, email string) (bool, error) {
	account, err := uc.accountRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrAccountNotFound) {
			return false, nil
		}
		return false, err
	}

	return account != nil, nil
}
