package credential

import (
	"auth-service/internal/model"
	"auth-service/internal/usecase/port"
	"auth-service/pkg/hasher"
	"context"
	"errors"
)

var (
	ErrAccountDisabled    = errors.New("account is disabled")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type CredentialVerifier struct {
	hasher      hasher.Hasher
	accountRepo port.AccountRepo
}

func NewCredentialVerifier(hasher hasher.Hasher, accountRepo port.AccountRepo) *CredentialVerifier {
	return &CredentialVerifier{
		hasher:      hasher,
		accountRepo: accountRepo,
	}
}

func (v *CredentialVerifier) Verify(ctx context.Context, email, password string) (*model.Account, error) {
	account, err := v.accountRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !account.IsActive {
		return nil, ErrAccountDisabled
	}

	match, err := v.hasher.ComparePasswordAndHash(password, account.PasswordHash)
	if err != nil || !match {
		return nil, ErrInvalidCredentials
	}

	return account, nil
}

func (v *CredentialVerifier) verify(ctx context.Context, email, password string) (*model.Account, error) {
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

func (uc *CredentialVerifier) checkAccountActive(account *model.Account) error {
	if !account.IsActive {
		return ErrAccountDisabled
	}

	return nil
}

func (v *CredentialVerifier) verifyPassword(plain, hashed string) error {
	match, err := v.hasher.ComparePasswordAndHash(plain, hashed)
	if err != nil {
		return err
	}
	if !match {
		return ErrInvalidCredentials
	}

	return nil
}
