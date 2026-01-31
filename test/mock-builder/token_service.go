package mockbuilder

import (
	"auth-service/internal/apperrs"
	"auth-service/internal/domain/accountmodel"
	"auth-service/test/fakes"
	"auth-service/test/mocks"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
)

var (
	ErrVerifyRefreshTokenFailed = errors.New("failed when verified refresh token")
)

type mockTokenService struct {
	inst *mocks.TokenService
}

func (b *mockTokenService) GetInstance() *mocks.TokenService {
	return b.inst
}

func newMockTokenService(t *testing.T) *mockTokenService {
	return &mockTokenService{
		inst: mocks.NewTokenService(t),
	}
}

func (b *mockTokenService) SignPairs_Failed(account *accountmodel.Account) {
	b.inst.EXPECT().
		SignPairs(mock.Anything, mock.Anything, account).
		Return(nil, apperrs.ErrInternal)
}

func (b *mockTokenService) SignPairs_Success(account *accountmodel.Account) {
	b.inst.EXPECT().
		SignPairs(mock.Anything, mock.Anything, account).
		Return(fakes.FakeTokenPairs(), nil)
}

func (b *mockTokenService) VerifyRefreshToken_Failed(ctx context.Context) {
	b.inst.EXPECT().VerifyRefreshToken(ctx, mock.Anything).Return(nil, ErrVerifyRefreshTokenFailed)
}

func (b *mockTokenService) VerifyRefreshToken_Success(ctx context.Context) {
	b.inst.EXPECT().VerifyRefreshToken(ctx, mock.AnythingOfType("string")).Return(fakes.FakeTokenClaims(), nil)
}

func (b *mockTokenService) Sign_Error(ctx context.Context) {
	b.inst.EXPECT().
		Sign(mock.Anything).
		Return("", ErrInternalDB).Once()
}

func (b *mockTokenService) Sign_OK(ctx context.Context) {
	b.inst.EXPECT().
		Sign(mock.Anything).
		Return(fakes.FakeTokenPairs().AccessToken, nil).Once()
}
