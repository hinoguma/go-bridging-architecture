package setup

import (
	"app/pkg/applogic/usecase"
	"context"
	"sync"
)

type UseCaseRegistry struct {
	sync.Mutex
	domainRepositoryRegistry *DomainRepositoryRegistry
	domainServiceRegistry    *DomainServiceRegistry

	depositUseCase           usecase.DepositUseCase
	bankAccountAuthMWUseCase usecase.BankAccountAuthMWUseCase
}

func (registry *UseCaseRegistry) Initialize(
	ctx context.Context,
	domainRepositoryRegistry *DomainRepositoryRegistry,
	domainServiceRegistry *DomainServiceRegistry,
) {
	registry.Lock()
	registry.domainServiceRegistry = domainServiceRegistry
	registry.domainRepositoryRegistry = domainRepositoryRegistry
	registry.depositUseCase = usecase.NewDepositUseCase(
		registry.domainServiceRegistry.DepositService(),
		registry.domainRepositoryRegistry.DBTransactionManager(),
	)
	registry.bankAccountAuthMWUseCase = usecase.NewBankAccountAuthMWUseCase(
		registry.domainServiceRegistry.BankAccountAuthService(),
	)
	registry.Unlock()
}

func (registry *UseCaseRegistry) DepositUseCase() usecase.DepositUseCase {
	if registry.depositUseCase == nil {
		registry.Lock()
		registry.depositUseCase = usecase.NewDepositUseCase(
			registry.domainServiceRegistry.DepositService(),
			registry.domainRepositoryRegistry.DBTransactionManager(),
		)
		registry.Unlock()
	}
	return registry.depositUseCase
}

func (registry *UseCaseRegistry) BankAccountAuthMWUseCase() usecase.BankAccountAuthMWUseCase {
	if registry.bankAccountAuthMWUseCase == nil {
		registry.Lock()
		registry.bankAccountAuthMWUseCase = usecase.NewBankAccountAuthMWUseCase(
			registry.domainServiceRegistry.BankAccountAuthService(),
		)
		registry.Unlock()
	}
	return registry.bankAccountAuthMWUseCase
}

func GetUseCaseRegistry() *UseCaseRegistry {
	return &useCaseRegistry
}

var useCaseRegistry UseCaseRegistry = UseCaseRegistry{}
