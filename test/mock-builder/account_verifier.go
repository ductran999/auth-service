package mockbuilder

import (
	"testing"

	"auth-service/internal/model"
	"auth-service/test/mocks"

	"github.com/stretchr/testify/mock"
)

type mockAccountVerifierBuilder struct {
	inst *mocks.AccountVerifier
}

func (b *mockAccountVerifierBuilder) GetInstance() *mocks.AccountVerifier {
	return b.inst
}

func newMockAccountVerifierBuilder(t *testing.T) *mockAccountVerifierBuilder {
	return &mockAccountVerifierBuilder{
		inst: mocks.NewAccountVerifier(t),
	}
}

func (b *mockAccountVerifierBuilder) VerifySuccess() {
	account := &model.Account{
		ID:       FakeAccountID,
		Email:    FakeEmail,
		IsActive: true,
	}

	b.inst.EXPECT().
		Verify(mock.Anything, mock.Anything, mock.Anything).
		Return(account, nil)
}

func (b *mockAccountVerifierBuilder) VerifyFailed(err error) {
	b.inst.EXPECT().
		Verify(mock.Anything, mock.Anything, mock.Anything).
		Return(nil, err)
}
