package setup

import (
	"app/internal/applogic/usecase"
	"context"
	"sync"
)

type UseCaseRegistry struct {
	sync.Mutex
	domainRepositoryRegistry *DomainRepositoryRegistry
	domainServiceRegistry    *DomainServiceRegistry

	bankAccountAuthMWUseCase usecase.BankAccountAuthMWUseCase
	depositUseCase           usecase.DepositUseCase
	withdrawUseCase          usecase.WithdrawUseCase
}

func (registry *UseCaseRegistry) Initialize(
	ctx context.Context,
	domainRepositoryRegistry *DomainRepositoryRegistry,
	domainServiceRegistry *DomainServiceRegistry,
) {
	registry.Lock()
	registry.domainServiceRegistry = domainServiceRegistry
	registry.domainRepositoryRegistry = domainRepositoryRegistry
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

func (registry *UseCaseRegistry) WithdrawUseCase() usecase.WithdrawUseCase {
	if registry.withdrawUseCase == nil {
		registry.Lock()
		registry.withdrawUseCase = usecase.NewWithdrawUseCase(
			registry.domainServiceRegistry.WithdrawService(),
			registry.domainRepositoryRegistry.DBTransactionManager(),
		)
		registry.Unlock()
	}
	return registry.withdrawUseCase
}

func GetUseCaseRegistry() *UseCaseRegistry {
	return &useCaseRegistry
}

var useCaseRegistry UseCaseRegistry = UseCaseRegistry{}
