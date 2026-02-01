package account

import (
	"auth-service/internal/apperrs"
	"auth-service/internal/biz/ports/repositories"
	"auth-service/internal/biz/ports/security"
	"context"
	"errors"
)

var (
	ErrAccountNotFound        = errors.New("account not found")
	ErrNewPasswordMustChanged = errors.New("new password must e different")
)

type ChangePasswordInput struct {
	AccountID   string
	OldPassword string
	NewPassword string
}

type changePasswordUsecase struct {
	hasher      security.PasswordHasher
	accountRepo repositories.AccountRepo
}

func (uc *changePasswordUsecase) ChangePassword(ctx context.Context, input ChangePasswordInput) error {
	account, err := uc.accountRepo.FindByID(ctx, input.AccountID)
	if err != nil {
		if errors.Is(err, ErrAccountNotFound) {
			return apperrs.NotFound(ErrAccountNotFound)
		}
		return apperrs.Internal(err)
	}

	if err = uc.validatePassword(input.OldPassword, account.PasswordHash); err != nil {
		return err
	}

	hashedPassword, err := uc.hashIfChanged(input.OldPassword, input.NewPassword)
	if err != nil {
		return err
	}

	err = uc.accountRepo.UpdatePasswordHash(ctx, account.ID.String(), hashedPassword)
	if err != nil {
		return apperrs.Internal(err)
	}

	return nil
}

func (uc *changePasswordUsecase) validatePassword(password, hashed string) error {
	match, err := uc.hasher.ComparePasswordAndHash(password, hashed)
	if err != nil {
		return apperrs.Internal(err)
	}
	if !match {
		return apperrs.Unauthorized(ErrPasswordMismatch)
	}

	return nil
}

func (uc *changePasswordUsecase) hashIfChanged(oldPassword, newPassword string) (string, error) {
	if oldPassword == newPassword {
		return "", apperrs.InvalidInput(ErrNewPasswordMustChanged)
	}

	passHash, err := uc.hasher.Hash(newPassword)
	if err != nil {
		return "", apperrs.Internal(err)

	}

	return passHash, nil
}
