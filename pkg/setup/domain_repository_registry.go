package setup

import (
	"app/pkg/applogic/domain/repository"
	"app/pkg/driven_infra_port"
	"context"
	"sync"
)

type DomainRepositoryRegistry struct {
	sync.Mutex
	drivenInfraRegistry         *DrivenInfraRegistry
	dbTransactionManager        repository.DBTransactionManager
	bankAccountRepository       repository.BankAccountRepository
	transactionRecordRepository repository.TransactionRecordRepository
	bankAccountAuthenticator    repository.BankAccountAuthenticator
}

func (registry *DomainRepositoryRegistry) Initialize(ctx context.Context, infraRegistry *DrivenInfraRegistry) {
	registry.Lock()
	registry.drivenInfraRegistry = infraRegistry
	registry.dbTransactionManager = driven_infra_port.NewDBTransactionManager(
		registry.drivenInfraRegistry.TransactionManager(),
	)
	registry.bankAccountRepository = driven_infra_port.NewBankAccountRepositorySQL(
		registry.drivenInfraRegistry.SQLClient(),
	)
	registry.transactionRecordRepository = driven_infra_port.NewTransactionRecordRepositorySQL(
		registry.drivenInfraRegistry.SQLClient(),
	)
	registry.bankAccountAuthenticator = driven_infra_port.NewBankAccountAuthenticator(
		registry.drivenInfraRegistry.CognitoClient(),
	)
	registry.Unlock()
}

func (registry *DomainRepositoryRegistry) DBTransactionManager() repository.DBTransactionManager {
	if registry.dbTransactionManager == nil {
		registry.Lock()
		registry.dbTransactionManager = driven_infra_port.NewDBTransactionManager(
			registry.drivenInfraRegistry.TransactionManager(),
		)
		registry.Unlock()
	}
	return registry.dbTransactionManager
}

func (registry *DomainRepositoryRegistry) BankAccountRepository() repository.BankAccountRepository {
	if registry.bankAccountRepository == nil {
		registry.Lock()
		registry.bankAccountRepository = driven_infra_port.NewBankAccountRepositoryDynamoDB(
			registry.drivenInfraRegistry.DynamoDBClient(),
		)
		registry.Unlock()
	}
	return registry.bankAccountRepository
}

func (registry *DomainRepositoryRegistry) TransactionRecordRepository() repository.TransactionRecordRepository {
	if registry.transactionRecordRepository == nil {
		registry.Lock()
		registry.transactionRecordRepository = driven_infra_port.NewTransactionRecordRepositoryDynamoDB(
			registry.drivenInfraRegistry.DynamoDBClient(),
		)
		registry.Unlock()
	}
	return registry.transactionRecordRepository
}

func (registry *DomainRepositoryRegistry) BankAccountAuthenticator() repository.BankAccountAuthenticator {
	if registry.bankAccountAuthenticator == nil {
		registry.Lock()
		registry.bankAccountAuthenticator = driven_infra_port.NewBankAccountAuthenticator(
			registry.drivenInfraRegistry.CognitoClient(),
		)
		registry.Unlock()
	}
	return registry.bankAccountAuthenticator
}

func GetDomainRepositoryRegistry() *DomainRepositoryRegistry {
	return &domainRepositoryRegistry
}

var domainRepositoryRegistry DomainRepositoryRegistry = DomainRepositoryRegistry{}
