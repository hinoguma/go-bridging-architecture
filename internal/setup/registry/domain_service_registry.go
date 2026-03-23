package registry

import (
	"app/internal/applogic/domain/service"
	"context"
	"sync"
)

type DomainServiceRegistry struct {
	sync.Mutex
	domainRepositoryRegistry *DomainRepositoryRegistry

	depositService         service.DepositService
	withdrawService        service.WithdrawService
	bankAccountAuthService service.BankAccountAuthService
}

func NewDomainServiceRegistry(
	domainRepositoryRegistry *DomainRepositoryRegistry,
) DomainServiceRegistry {
	return DomainServiceRegistry{
		domainRepositoryRegistry: domainRepositoryRegistry,
	}
}

func (registry *DomainServiceRegistry) Initialize(
	ctx context.Context,
	domainRepositoryRegistry *DomainRepositoryRegistry,
) {
	registry.Lock()
	registry.domainRepositoryRegistry = domainRepositoryRegistry
	registry.Unlock()
}

func (registry *DomainServiceRegistry) DepositService() service.DepositService {
	if registry.depositService == nil {
		registry.Lock()
		registry.depositService = service.NewDepositService(
			registry.domainRepositoryRegistry.DBTransactionManager(),
			registry.domainRepositoryRegistry.BankAccountRepository(),
			registry.domainRepositoryRegistry.TransactionRecordRepository(),
		)
		registry.Unlock()
	}
	return registry.depositService
}

func (registry *DomainServiceRegistry) WithdrawService() service.WithdrawService {
	if registry.withdrawService == nil {
		registry.Lock()
		registry.withdrawService = service.NewWithdrawService(
			registry.domainRepositoryRegistry.DBTransactionManager(),
			registry.domainRepositoryRegistry.BankAccountRepository(),
			registry.domainRepositoryRegistry.TransactionRecordRepository(),
		)
		registry.Unlock()
	}
	return registry.withdrawService
}

func (registry *DomainServiceRegistry) BankAccountAuthService() service.BankAccountAuthService {
	if registry.bankAccountAuthService == nil {
		registry.Lock()
		registry.bankAccountAuthService = service.NewBankAccountAuthService(
			registry.domainRepositoryRegistry.BankAccountRepository(),
			registry.domainRepositoryRegistry.BankAccountAuthenticator(),
		)
		registry.Unlock()
	}
	return registry.bankAccountAuthService
}

var domainServiceRegistryInstance DomainServiceRegistry = DomainServiceRegistry{}

func GetDomainServiceRegistry() *DomainServiceRegistry {
	return &domainServiceRegistryInstance
}
