package account

import (
	"auth-service/internal/usecase/port"
	"auth-service/pkg/hasher"
)

type accountUsecase struct {
	*registerUsecase
	*changePasswordUsecase
}

func NewAccountUseCase(hasher hasher.Hasher, accountRepo port.AccountRepo) port.AccountUsecase {
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
