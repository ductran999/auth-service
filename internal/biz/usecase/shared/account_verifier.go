package shared

// import (
// 	"context"

// 	"auth-service/internal/model"
// 	"auth-service/internal/usecase/port"
// 	"auth-service/pkg/hasher"
// )

// type AccountVerifier interface {
// 	Verify(ctx context.Context, email, password string) (*model.Account, error)
// }

// type accountVerifier struct {
// 	hasher      hasher.Hasher
// 	accountRepo port.AccountRepo
// }

// func NewAccountVerifier(hasher hasher.Hasher, accountRepo port.AccountRepo) AccountVerifier {
// 	return &accountVerifier{
// 		hasher:      hasher,
// 		accountRepo: accountRepo,
// 	}
// }

// func (v *accountVerifier) Verify(ctx context.Context, email, password string) (*model.Account, error) {
// 	account, err := v.accountRepo.FindByEmail(ctx, email)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err := v.checkAccountActive(account); err != nil {
// 		return nil, err
// 	}

// 	if err := v.verifyPassword(password, account.PasswordHash); err != nil {
// 		return nil, err
// 	}

// 	return account, nil
// }

// func (uc *accountVerifier) checkAccountActive(account *model.Account) error {
// 	if !account.IsActive {
// 		return ErrAccountDisabled
// 	}

// 	return nil
// }

// func (uc *accountVerifier) verifyPassword(plain, hashed string) error {
// 	_, err := uc.hasher.ComparePasswordAndHash(plain, hashed)
// 	if err != nil {
// 		return err
// 	}
// 	// if !match {
// 	// 	return sessionUC.ErrInvalidCredentials
// 	// }

// 	return nil
// }
