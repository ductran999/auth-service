package mockbuilder

import "testing"

type BuilderContainer struct {
	AccountRepoBuilder *mockAccountRepoBuilder
	SessionRepoBuilder *mockSessionRepoBuilder
	HasherBuilder      *mockHasherBuilder
	// CacheBuilder       *mockCacheBuilder
	// TokenSigner     *mockSignerBuilder
}

func NewBuilderContainer(t *testing.T) *BuilderContainer {
	return &BuilderContainer{
		AccountRepoBuilder: newMockAccountRepoBuilder(t),
		SessionRepoBuilder: newMockSessionRepoBuilder(t),
		HasherBuilder:      newMockHasherBuilder(t),
		// CacheBuilder:       newMockCacheBuilder(t),
		// TokenSigner:     NewMockSignerBuilder(t),
	}
}

type UsecaseBuilderContainer struct {
	AccountUC *mockAccountUsecase
	SessionUC *mockSessionUsecase
	// AuthJwtUC     *mockAuthJWTUsecase
	// AuthSessionUC *mockAuthSessionUsecase
}

func NewUsecaseBuilderContainer(t *testing.T) *UsecaseBuilderContainer {
	return &UsecaseBuilderContainer{
		AccountUC: newMockAccountUsecase(t),
		SessionUC: newMockSessionUsecase(t),
		// AuthJwtUC:     newMockAuthJWTUsecase(t),
		// AuthSessionUC: newMockAuthSessionUsecase(t),
	}
}
