package account

import (
	"context"
	"errors"

	"auth-service/internal/biz/ports/repositories"
	accountUC "auth-service/internal/biz/usecase/account"
	"auth-service/internal/domain/accountmodel"

	"gorm.io/gorm"
)

type accountRepo struct {
	db *gorm.DB
}

func NewAccountRepo(db *gorm.DB) repositories.AccountRepo {
	return &accountRepo{db: db}
}

// Create inserts a new account record into the database.
func (r *accountRepo) Create(ctx context.Context, account *accountmodel.Account) error {
	return r.db.WithContext(ctx).Create(account).Error
}

// FindByEmail looks up an account by its email address.
func (r *accountRepo) FindByEmail(ctx context.Context, email string) (*accountmodel.Account, error) {
	var account accountmodel.Account

	err := r.db.WithContext(ctx).First(&account, "email = ?", email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, accountUC.ErrAccountNotFound
		}
		return nil, err
	}

	return &account, nil
}

func (r *accountRepo) FindByID(ctx context.Context, id string) (*accountmodel.Account, error) {
	var account accountmodel.Account

	err := r.db.WithContext(ctx).First(&account, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, accountUC.ErrAccountNotFound
		}
		return nil, err
	}

	return &account, nil
}

func (r *accountRepo) UpdatePasswordHash(ctx context.Context, id, passwordHash string) error {
	return r.db.WithContext(ctx).
		Model(&accountmodel.Account{}).
		Where("id = ?", id).
		Update("password_hash", passwordHash).
		Error
}
