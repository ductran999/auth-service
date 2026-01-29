package credential

import (
	"auth-service/internal/biz/ports/repositories"
	"auth-service/internal/biz/ports/security"
	"auth-service/internal/domain/accountmodel"
	"context"
	"errors"
)

var (
	ErrAccountDisabled    = errors.New("account is disabled")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type CredentialVerifier struct {
	hasher      security.PasswordHasher
	accountRepo repositories.AccountRepo
}

func NewCredentialVerifier(hasher security.PasswordHasher, accountRepo repositories.AccountRepo) *CredentialVerifier {
	return &CredentialVerifier{
		hasher:      hasher,
		accountRepo: accountRepo,
	}
}

func (v *CredentialVerifier) Verify(ctx context.Context, email, password string) (*accountmodel.Account, error) {
	account, err := v.accountRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	if !account.IsActive {
		return nil, ErrAccountDisabled
	}

	match, err := v.hasher.ComparePasswordAndHash(password, account.PasswordHash)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, ErrInvalidCredentials
	}

	return account, nil
}
