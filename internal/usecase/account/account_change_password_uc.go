package account

import (
	"auth-service/internal/apperrs"
	"auth-service/internal/usecase/dto"
	"auth-service/internal/usecase/port"
	"auth-service/pkg/hasher"
	"context"
	"errors"
)

type changePasswordUsecase struct {
	hasher      hasher.Hasher
	accountRepo port.AccountRepo
}

func (uc *changePasswordUsecase) ChangePassword(ctx context.Context, input dto.ChangePasswordInput) error {
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

	passHash, err := uc.hasher.HashPassword(newPassword)
	if err != nil {
		return "", apperrs.Internal(err)

	}

	return passHash, nil
}
