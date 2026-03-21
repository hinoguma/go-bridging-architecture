package setup

import (
	"app/internal/applogic/domain/repository"
	"app/internal/applogic_appinfra_connector"
	"context"
	"sync"
)

type DomainRepositoryRegistry struct {
	sync.Mutex
	appInfraRegistry            *AppInfraRegistry
	dbTransactionManager        repository.DBTransactionManager
	bankAccountRepository       repository.BankAccountRepository
	transactionRecordRepository repository.TransactionRecordRepository
	bankAccountAuthenticator    repository.BankAccountAuthenticator
}

func (registry *DomainRepositoryRegistry) Initialize(ctx context.Context, infraRegistry *AppInfraRegistry) {
	registry.Lock()
	registry.appInfraRegistry = infraRegistry
	registry.dbTransactionManager = applogic_appinfra_connector.NewDBTransactionManager(
		registry.appInfraRegistry.TransactionManager(),
	)
	registry.bankAccountRepository = applogic_appinfra_connector.NewBankAccountRepositoryPostgresConnector(
		registry.appInfraRegistry.SQLClient(),
	)
	registry.transactionRecordRepository = applogic_appinfra_connector.NewTransactionRecordRepositoryPostgresConnector(
		registry.appInfraRegistry.SQLClient(),
	)
	registry.bankAccountAuthenticator = applogic_appinfra_connector.NewBankAccountAuthenticator(
		registry.appInfraRegistry.CognitoClient(),
	)
	registry.Unlock()
}

func (registry *DomainRepositoryRegistry) DBTransactionManager() repository.DBTransactionManager {
	if registry.dbTransactionManager == nil {
		registry.Lock()
		registry.dbTransactionManager = applogic_appinfra_connector.NewDBTransactionManager(
			registry.appInfraRegistry.TransactionManager(),
		)
		registry.Unlock()
	}
	return registry.dbTransactionManager
}

func (registry *DomainRepositoryRegistry) BankAccountRepository() repository.BankAccountRepository {
	if registry.bankAccountRepository == nil {
		registry.Lock()
		registry.bankAccountRepository = applogic_appinfra_connector.NewBankAccountRepositoryDynamoDBConnector(
			registry.appInfraRegistry.DynamoDBClient(),
		)
		registry.Unlock()
	}
	return registry.bankAccountRepository
}

func (registry *DomainRepositoryRegistry) TransactionRecordRepository() repository.TransactionRecordRepository {
	if registry.transactionRecordRepository == nil {
		registry.Lock()
		registry.transactionRecordRepository = applogic_appinfra_connector.NewTransactionRecordRepositoryDynamoDBConnector(
			registry.appInfraRegistry.DynamoDBClient(),
		)
		registry.Unlock()
	}
	return registry.transactionRecordRepository
}

func (registry *DomainRepositoryRegistry) BankAccountAuthenticator() repository.BankAccountAuthenticator {
	if registry.bankAccountAuthenticator == nil {
		registry.Lock()
		registry.bankAccountAuthenticator = applogic_appinfra_connector.NewBankAccountAuthenticator(
			registry.appInfraRegistry.CognitoClient(),
		)
		registry.Unlock()
	}
	return registry.bankAccountAuthenticator
}

func GetDomainRepositoryRegistry() *DomainRepositoryRegistry {
	return &domainRepositoryRegistry
}

var domainRepositoryRegistry DomainRepositoryRegistry = DomainRepositoryRegistry{}
