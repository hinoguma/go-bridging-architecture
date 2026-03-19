package setup

import (
	"app/internal/applogic/domain/repository"
	"app/internal/applogic_appinfra_bridge"
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
	registry.dbTransactionManager = applogic_appinfra_bridge.NewDBTransactionManager(
		registry.drivenInfraRegistry.TransactionManager(),
	)
	registry.bankAccountRepository = applogic_appinfra_bridge.NewBankAccountRepositorySQL(
		registry.drivenInfraRegistry.SQLClient(),
	)
	registry.transactionRecordRepository = applogic_appinfra_bridge.NewTransactionRecordRepositorySQL(
		registry.drivenInfraRegistry.SQLClient(),
	)
	registry.bankAccountAuthenticator = applogic_appinfra_bridge.NewBankAccountAuthenticator(
		registry.drivenInfraRegistry.CognitoClient(),
	)
	registry.Unlock()
}

func (registry *DomainRepositoryRegistry) DBTransactionManager() repository.DBTransactionManager {
	if registry.dbTransactionManager == nil {
		registry.Lock()
		registry.dbTransactionManager = applogic_appinfra_bridge.NewDBTransactionManager(
			registry.drivenInfraRegistry.TransactionManager(),
		)
		registry.Unlock()
	}
	return registry.dbTransactionManager
}

func (registry *DomainRepositoryRegistry) BankAccountRepository() repository.BankAccountRepository {
	if registry.bankAccountRepository == nil {
		registry.Lock()
		registry.bankAccountRepository = applogic_appinfra_bridge.NewBankAccountRepositoryDynamoDB(
			registry.drivenInfraRegistry.DynamoDBClient(),
		)
		registry.Unlock()
	}
	return registry.bankAccountRepository
}

func (registry *DomainRepositoryRegistry) TransactionRecordRepository() repository.TransactionRecordRepository {
	if registry.transactionRecordRepository == nil {
		registry.Lock()
		registry.transactionRecordRepository = applogic_appinfra_bridge.NewTransactionRecordRepositoryDynamoDB(
			registry.drivenInfraRegistry.DynamoDBClient(),
		)
		registry.Unlock()
	}
	return registry.transactionRecordRepository
}

func (registry *DomainRepositoryRegistry) BankAccountAuthenticator() repository.BankAccountAuthenticator {
	if registry.bankAccountAuthenticator == nil {
		registry.Lock()
		registry.bankAccountAuthenticator = applogic_appinfra_bridge.NewBankAccountAuthenticator(
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
