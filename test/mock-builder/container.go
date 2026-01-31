package mockbuilder

import "testing"

type BuilderContainer struct {
	AccountRepoBuilder *mockAccountRepoBuilder
	SessionRepoBuilder *mockSessionRepoBuilder
	HasherBuilder      *mockHasherBuilder

	TokenService *mockTokenService
	TokenStore   *mockTokenStore

	SessionStore *mockSessionStore
}

func NewBuilderContainer(t *testing.T) *BuilderContainer {
	return &BuilderContainer{
		AccountRepoBuilder: newMockAccountRepoBuilder(t),
		SessionRepoBuilder: newMockSessionRepoBuilder(t),
		HasherBuilder:      newMockHasherBuilder(t),
		TokenService:       newMockTokenService(t),
		TokenStore:         newMockTokenStore(t),
		SessionStore:       newMockSessionStore(t),
	}
}

type UsecaseBuilderContainer struct {
	AccountUC     *mockAccountUsecase
	AuthSessionUC *mockAuthSessionUsecase
}

func NewUsecaseBuilderContainer(t *testing.T) *UsecaseBuilderContainer {
	return &UsecaseBuilderContainer{
		AccountUC:     newMockAccountUsecase(t),
		AuthSessionUC: newMockAuthSessionUsecase(t),
	}
}
