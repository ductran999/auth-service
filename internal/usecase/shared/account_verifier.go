package shared

import (
	"context"

	"auth-service/internal/errs"
	"auth-service/internal/model"
	"auth-service/internal/usecase/port"
	"auth-service/pkg/hasher"
)

type AccountVerifier interface {
	Verify(ctx context.Context, email, password string) (*model.Account, error)
}

type accountVerifier struct {
	hasher      hasher.Hasher
	accountRepo port.AccountRepo
}

func NewAccountVerifier(hasher hasher.Hasher, accountRepo port.AccountRepo) AccountVerifier {
	return &accountVerifier{
		hasher:      hasher,
		accountRepo: accountRepo,
	}
}

func (v *accountVerifier) Verify(ctx context.Context, email, password string) (*model.Account, error) {
	account, err := v.accountRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err := v.checkAccountActive(account); err != nil {
		return nil, err
	}

	if err := v.verifyPassword(password, account.PasswordHash); err != nil {
		return nil, err
	}

	return account, nil
}

func (uc *accountVerifier) checkAccountActive(account *model.Account) error {
	if !account.IsActive {
		return errs.ErrAccountDisabled
	}
	return nil
}

func (uc *accountVerifier) verifyPassword(plain, hashed string) error {
	match, err := uc.hasher.ComparePasswordAndHash(plain, hashed)
	if err != nil {
		return err
	}
	if !match {
		return errs.ErrInvalidCredentials
	}
	return nil
}
