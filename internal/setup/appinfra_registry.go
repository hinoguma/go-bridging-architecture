package setup

import (
	"app/internal/appinfra/aws"
	"app/internal/appinfra/postgres"
	"app/internal/crosscutting/errors"
	"context"
	"database/sql"
	"sync"
)

type AppInfraRegistry struct {
	sync.Mutex
	dynamoDBClient     driven_infra_aws.DynamoDBClient
	cognitoClient      driven_infra_aws.CognitoClient
	sqlClient          postgres.SQLClient
	transactionManager postgres.TransactionManagerIF
}

func (registry *AppInfraRegistry) Initialize(
	ctx context.Context,
	db *sql.DB,
) error {
	var err error
	registry.Lock()
	registry.dynamoDBClient, err = driven_infra_aws.NewDynamoDBClient(ctx)
	if err != nil {
		registry.Unlock()
		return errors.Lift(err)
	}
	registry.cognitoClient, err = driven_infra_aws.NewCognitoClient(ctx)
	if err != nil {
		registry.Unlock()
		return errors.Lift(err)
	}
	registry.sqlClient = postgres.NewPostgreSQLClient(db)
	registry.transactionManager = postgres.NewTransactionManager(registry.sqlClient.GetDB())
	registry.Unlock()
	return nil
}

func (registry *AppInfraRegistry) DynamoDBClient() driven_infra_aws.DynamoDBClient {
	return registry.dynamoDBClient
}

func (registry *AppInfraRegistry) CognitoClient() driven_infra_aws.CognitoClient {
	return registry.cognitoClient
}

func (registry *AppInfraRegistry) SQLClient() postgres.SQLClient {
	return registry.sqlClient
}

func (registry *AppInfraRegistry) TransactionManager() postgres.TransactionManagerIF {
	return registry.transactionManager
}

var appInfraRegistry = AppInfraRegistry{}

func GetAppInfraRegistry() *AppInfraRegistry {
	return &appInfraRegistry
}
